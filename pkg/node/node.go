package node

import (
	"cmp/common"
	"cmp/model/k8s"
	"cmp/pkg/evict"
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type K8sNodeList struct {
	Nodes []Node `json:"nodes"`
}
type Node struct {
	ObjectMeta         k8s.ObjectMeta             `json:"objectMeta"`
	TypeMeta           k8s.TypeMeta               `json:"typeMeta"`           //通过label判断role类型
	Ready              v1.ConditionStatus         `json:"ready"`              //节点状态
	Unschedulable      k8s.Unschedulable          `json:"unschedulable"`      //是否可以调度
	NodeIP             k8s.NodeIP                 `json:"nodeIP"`             //node IP
	AllocatedResources k8s.NodeAllocatedResources `json:"allocatedResources"` //节点资源
	NodeInfo           v1.NodeSystemInfo          `json:"nodeInfo"`
}

func GetNodeList(client *kubernetes.Clientset) (*K8sNodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get nodes from cluster failed: %v", err)
	}
	return toNodeList(client, nodes.Items), nil
}

func toNodeList(client *kubernetes.Clientset, nodes []v1.Node) *K8sNodeList {

	nodeList := &K8sNodeList{
		Nodes: make([]Node, 0), // make初始化node信息
	}

	for _, node := range nodes {
		// 根据Node名称去获取节点上面的pod，过滤时排除pod为 Succeeded, Failed 。返回pods
		pods, err := getNodePods(client, node)
		if err != nil {
			common.Log.Error(fmt.Sprintf("Couldn't get pods of %s node: %s\n", node.Name, err))
		}

		// 调用toNode方法获取 node节点的计算资源
		nodeList.Nodes = append(nodeList.Nodes, toNode(node, pods, getNodeRole(node)))
	}

	return nodeList
}

func toNode(node v1.Node, pods *v1.PodList, role string) Node {
	// 获取cpu和内存的reqs, limits使用
	allocatedResources, err := getNodeAllocatedResources(node, pods)
	if err != nil {
		common.Log.Error(fmt.Sprintf("Couldn't get allocated resources of %s node: %s\n", node.Name, err))
	}

	return Node{
		ObjectMeta:         k8s.NewObjectMeta(node.ObjectMeta),
		TypeMeta:           k8s.NewTypeMeta(k8s.ResourceKind(role)),
		Ready:              getNodeConditionStatus(node, v1.NodeReady),
		NodeIP:             k8s.NodeIP(getNodeIP(node)),
		Unschedulable:      k8s.Unschedulable(node.Spec.Unschedulable),
		AllocatedResources: allocatedResources,
		NodeInfo:           node.Status.NodeInfo,
	}
}

func getNodeConditionStatus(node v1.Node, conditionType v1.NodeConditionType) v1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return v1.ConditionUnknown
}

func getNodeIP(node v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}

func getNodeRole(node v1.Node) string {
	var role string
	if _, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]; ok {
		role = "Master"
	} else {
		role = "Worker"
	}
	return role
}

func K8sNodeUnschedulable(client *kubernetes.Clientset, nodeName string, unschdulable bool) (bool, error) {
	//设置节点是否可调度
	common.Log.Info(fmt.Sprintf("设置Node节点:%v 是否可调度:%v", nodeName, unschdulable))
	node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		common.Log.Error(fmt.Sprintf("获取node失败: %v", err.Error()))
		return false, err
	}
	node.Spec.Unschedulable = unschdulable
	_, err = client.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		common.Log.Error(fmt.Sprintf("设置节点调度失败:%v", err.Error()))
		return false, err
	}
	return true, nil
}

// CordonNode 选择排空节点（同时设置为不可调度），在后续进行应用部署时，则Pod不会再调度到该节点，并且该节点上由DaemonSet控制的Pod不会被排空。
func CordonNode(client *kubernetes.Clientset, nodeName string) (bool, error) {
	_, err := K8sNodeUnschedulable(client, nodeName, true)
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

func RemoveNode(client *kubernetes.Clientset, nodeName string) (bool, error) {
	startTime := time.Now()
	common.Log.Info(fmt.Sprintf("移除Node节点：%v,异步任务已开始", nodeName))
	_, err := K8sNodeUnschedulable(client, nodeName, true)
	if err != nil {
		return false, err
	}
	err = evict.EvictsNodePods(client, nodeName)
	if err != nil {
		common.Log.Error(fmt.Sprintf("排空节点出现异常:%v", err.Error()))
		return false, nil
	}
	err2 := client.CoreV1().Nodes().Delete(context.TODO(), nodeName, metav1.DeleteOptions{})
	if err2 != nil {
		return false, nil
	}
	common.Log.Info(fmt.Sprintf("已将节点%v从集群中移除，异步任务已完成，任务耗时：%v", nodeName, time.Since(startTime)))
	return true, nil
}

// CollectionNodeUnschedule 批量设置node节点不可调度
func CollectionNodeUnschedule(client *kubernetes.Clientset, nodeName []string) error {
	if len(nodeName) <= 0 {
		return errors.New("节点名称不能为空")
	}
	common.Log.Info(fmt.Sprintf("批量设置Node节点：%v 为不可调度：true", nodeName))
	for _, v := range nodeName {
		node, err := client.CoreV1().Nodes().Get(context.TODO(), v, metav1.GetOptions{})
		if err != nil {
			common.Log.Error(fmt.Sprintf("获取node失败：%v", err.Error()))
			return err
		}
		node.Spec.Unschedulable = true
		_, err2 := client.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
		if err2 != nil {
			common.Log.Error(fmt.Sprintf("设置节点调度失败：%v", err2.Error()))
			return err
		}
	}
	common.Log.Info(fmt.Sprintf("已将所有Node节点:%v 设置为不可调度", nodeName))
	return nil
}

// CollectionCordonNode 批量排空node节点，不允许调度
func CollectionCordonNode(client *kubernetes.Clientset, nodeName []string) error {
	if len(nodeName) <= 0 {
		return errors.New("节点名称不能为空")
	}
	common.Log.Info(fmt.Sprintf("开始排空节点, 设置Node节点:%v  不可调度：true", nodeName))
	for _, v := range nodeName {
		node, err := client.CoreV1().Nodes().Get(context.TODO(), v, metav1.GetOptions{})
		if err != nil {
			common.Log.Error(fmt.Sprintf("获取节点失败: %v", err.Error()))
			return err
		}
		node.Spec.Unschedulable = true

		_, err2 := client.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

		if err2 != nil {
			common.Log.Error(fmt.Sprintf("设置节点调度失败：%v", err2.Error()))
			return err
		}
		_, cordonErr := CordonNode(client, v)
		if cordonErr != nil {
			return cordonErr
		}
	}
	common.Log.Info(fmt.Sprintf("已将所有Node节点:%v  设置为不可调度", nodeName))
	return nil
}
