package node

import (
	"cmp/common"
	"cmp/model/k8s"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeList struct {
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

func GetNodeList(client *kubernetes.Clientset) (*NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get nodes from cluster failed: %v", err)
	}
	return toNodeList(client, nodes.Items), nil
}

func toNodeList(client *kubernetes.Clientset, nodes []v1.Node) *NodeList {

	nodeList := &NodeList{
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
