package node

import (
	"cmp/common"
	"cmp/model/k8s"
	pkgcommon "cmp/pkg/common"
	"cmp/pkg/event"
	"cmp/pkg/evict"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

type nodeAllocatedResources struct {
	CPURequests            int64   `json:"cpuRequests"`
	CPURequestsFraction    float64 `json:"cpuRequestsFraction"`
	CPULimits              int64   `json:"cpuLimits"`
	CPULimitsFraction      float64 `json:"cpuLimitsFraction"`
	CPUCapacity            int64   `json:"cpuCapacity"`
	MemoryRequests         int64   `json:"memoryRequests"`
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`
	MemoryLimits           int64   `json:"memoryLimits"`
	MemoryLimitsFraction   float64 `json:"memoryLimitsFraction"`
	MemoryCapacity         int64   `json:"memoryCapacity"`
	AllocatedPods          int     `json:"allocatedPods"`
	PodCapacity            int64   `json:"podCapacity"`
	PodFraction            float64 `json:"podFraction"`
}

type nodeDetail struct {
	// Extends list item structure.
	Node `json:",inline"`
	// PodCIDR represents the pod IP range assigned to the node.
	PodCIDR string `json:"podCIDR"`
	// ID of the node assigned by the cloud provider.
	ProviderID string `json:"providerID"`
	// Unschedulable controls node schedulability of new pods. By default node is schedulable.
	Unschedulable bool `json:"unschedulable"`
	// Set of ids/uuids to uniquely identify the node.
	NodeInfo v1.NodeSystemInfo `json:"nodeInfo"`
	//// Conditions is an array of current node conditions.
	Conditions []pkgcommon.Condition `json:"conditions"`
	// Container images of the node.
	ContainerImages []string `json:"containerImages"`
	// PodListComponent contains information about pods belonging to this node.
	PodList v1.PodList `json:"podList"`
	// Events is list of events associated to the node.
	EventList v1.EventList `json:"eventList"`
	// Taints
	Taints []v1.Taint `json:"taints,omitempty"`
	// Addresses is a list of addresses reachable to the node. Queried from cloud provider, if available.
	Addresses []v1.NodeAddress   `json:"addresses,omitempty"`
	Ready     v1.ConditionStatus `json:"ready"`
	NodeIP    k8s.NodeIP         `json:"nodeIP"`
	UID       k8s.UID            `json:"uid"`
}

func GetNodeDetail(client *kubernetes.Clientset, name string) (*nodeDetail, error) {
	common.Log.Info(fmt.Sprintf("Getting details of %s node", name))
	node, err := client.CoreV1().Nodes().Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	pods, err := getNodePods(client, *node)
	if err != nil {
		return nil, err
	}
	eventList, err := event.GetNodeEvents(client, node.Name)
	if err != nil {
		return nil, err
	}
	allocatedResources, err := getNodeAllocatedResources(*node, pods)
	if err != nil {
		return nil, err
	}
	nodeDetails := toNodeDetail(*node, pods, eventList, nodeAllocatedResources(allocatedResources))
	return &nodeDetails, nil
}

func getNodePods(client *kubernetes.Clientset, node v1.Node) (*v1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name +
		",status.phase!=" + string(v1.PodSucceeded) +
		",status.phase!=" + string(v1.PodFailed))

	if err != nil {
		return nil, err
	}

	return client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), metaV1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
}

// 获取cpu和内存的reqs, limits使用
func getNodeAllocatedResources(node v1.Node, podList *v1.PodList) (k8s.NodeAllocatedResources, error) {
	reqs, limits := map[v1.ResourceName]resource.Quantity{}, map[v1.ResourceName]resource.Quantity{}
	// 遍历pod list获取pod cpu和内存的reqs, limits使用
	for _, pod := range podList.Items {
		podReqs, podLimits, err := PodRequestsAndLimits(&pod)
		if err != nil {
			return k8s.NodeAllocatedResources{}, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}

	cpuRequests, cpuLimits, memoryRequests, memoryLimits := reqs[v1.ResourceCPU], limits[v1.ResourceCPU], reqs[v1.ResourceMemory], limits[v1.ResourceMemory]

	var cpuRequestsFraction, cpuLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Allocatable.Cpu().MilliValue()); capacity > 0 {
		cpuRequestsFraction = float64(cpuRequests.MilliValue()) / capacity * 100
		cpuLimitsFraction = float64(cpuLimits.MilliValue()) / capacity * 100
	}

	var memoryRequestsFraction, memoryLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Allocatable.Memory().MilliValue()); capacity > 0 {
		memoryRequestsFraction = float64(memoryRequests.MilliValue()) / capacity * 100
		memoryLimitsFraction = float64(memoryLimits.MilliValue()) / capacity * 100
	}

	var podFraction float64 = 0
	var podCapacity = node.Status.Capacity.Pods().Value()
	if podCapacity > 0 {
		podFraction = float64(len(podList.Items)) / float64(podCapacity) * 100
	}

	return k8s.NodeAllocatedResources{
		CPURequests:            cpuRequests.MilliValue(),
		CPURequestsFraction:    cpuRequestsFraction,
		CPULimits:              cpuLimits.MilliValue(),
		CPULimitsFraction:      cpuLimitsFraction,
		CPUCapacity:            node.Status.Allocatable.Cpu().MilliValue(),
		MemoryRequests:         memoryRequests.Value(),
		MemoryRequestsFraction: memoryRequestsFraction,
		MemoryLimits:           memoryLimits.Value(),
		MemoryLimitsFraction:   memoryLimitsFraction,
		MemoryCapacity:         node.Status.Allocatable.Memory().Value(),
		AllocatedPods:          len(podList.Items),
		PodCapacity:            podCapacity,
		PodFraction:            podFraction,
	}, nil
}

// PodRequestsAndLimits 获取pod cpu和内存的reqs, limits使用
func PodRequestsAndLimits(pod *v1.Pod) (reqs, limits v1.ResourceList, err error) {
	reqs, limits = v1.ResourceList{}, v1.ResourceList{}
	//遍历pod中的容器
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// 初始化容器定义任何资源的最小值
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}

	// 将运行 pod 的开销添加到请求总和和非零限制
	if pod.Spec.Overhead != nil {
		addResourceList(reqs, pod.Spec.Overhead)

		for name, quantity := range pod.Spec.Overhead {
			if value, ok := limits[name]; ok && !value.IsZero() {
				value.Add(quantity)
				limits[name] = value
			}
		}
	}
	return
}
func addResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}
func maxResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = quantity.DeepCopy()
			}
		}
	}
}

func toNodeDetail(node v1.Node, pods *v1.PodList, eventList *v1.EventList, allocatedResources nodeAllocatedResources) nodeDetail {
	return nodeDetail{
		Node: Node{
			ObjectMeta:         k8s.NewObjectMeta(node.ObjectMeta),
			TypeMeta:           k8s.NewTypeMeta("node"),
			AllocatedResources: k8s.NodeAllocatedResources(allocatedResources),
		},
		ProviderID:      node.Spec.ProviderID,
		PodCIDR:         node.Spec.PodCIDR,
		Unschedulable:   node.Spec.Unschedulable,
		NodeInfo:        node.Status.NodeInfo,
		Conditions:      getNodeConditions(node),
		ContainerImages: getContainerImages(node),
		PodList:         *pods,
		EventList:       *eventList,
		Taints:          node.Spec.Taints,
		Addresses:       node.Status.Addresses,
		Ready:           getNodeConditionStatus(node, v1.NodeReady),
		NodeIP:          k8s.NodeIP(getNodeIP(node)),
		UID:             k8s.UID(node.UID),
	}
}

func NodeUnschedulable(client *kubernetes.Clientset, nodeName string, unschdulable bool) (bool, error) {
	//设置节点是否可调度
	common.Log.Info(fmt.Sprintf("设置Node节点:%v 不可调度:%v", nodeName, unschdulable))
	node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err != nil {
		common.Log.Error(fmt.Sprintf("get node err: %v", err.Error()))
		return false, err
	}
	node.Spec.Unschedulable = unschdulable
	_, err = client.CoreV1().Nodes().Update(context.TODO(), node, metaV1.UpdateOptions{})
	if err != nil {
		common.Log.Error(fmt.Sprintf("设置节点调度失败:%v", err.Error()))
		return false, err
	}
	return true, nil
}

func CordonNode(client *kubernetes.Clientset, nodeName string) (bool, error) {
	//排空节点
	_, err := NodeUnschedulable(client, nodeName, true)
	if err != nil {
		return false, nil
	}
	//驱逐节点上不在 kube-system 命名空间中的所有 pod
	err = evict.EvictsNodePods(client, nodeName)
	if err != nil {
		common.Log.Error(fmt.Sprintf("排空节点出现异常: %v", err.Error()))
		return false, err
	}
	return true, nil
}
