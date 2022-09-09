package cluster

import (
	"cmp/common"
	"cmp/model/k8s"
	"context"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetClusterStatus(c *kubernetes.Clientset) (string, error) {
	version, err := c.ServerVersion()
	if err != nil {
		common.Log.Error("get version from cluster failed", zap.Any("err: ", err))
		return "", err
	}
	return version.String(), nil
}

func GetClusterNumber(c *kubernetes.Clientset) (int, error) {
	number, err := c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(number.Items), nil
}

func GetClusterInfo(c *kubernetes.Clientset) *k8s.ClusterStatus {
	var node k8s.ClusterStatus

	nodesList, err := c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil
	}
	nodes := nodesList.Items

	totalCpu := float64(0)
	totalMemory := float64(0)
	usedCpu := float64(0)
	usedMemory := float64(0)
	readyNodes := 0
	unreadyNodes := 0

	for i := range nodes {
		conditions := nodes[i].Status.Conditions
		for i := range conditions {
			if conditions[i].Type == "Ready" {
				if conditions[i].Status == "True" {
					readyNodes += 1
				} else {
					unreadyNodes += 1
				}
			}
		}
		cpu := nodes[i].Status.Allocatable.Cpu().AsApproximateFloat64()
		totalCpu += cpu
		memory := nodes[i].Status.Allocatable.Memory().AsApproximateFloat64()
		totalMemory += memory
	}
	podsList, err := c.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil
	}
	pods := podsList.Items
	for i := range pods {
		for j := range pods[i].Spec.Containers {
			cpu := pods[i].Spec.Containers[j].Resources.Requests.Cpu().AsApproximateFloat64()
			usedCpu += cpu
			memory := pods[i].Spec.Containers[j].Resources.Requests.Memory().AsApproximateFloat64()
			usedMemory += memory

		}
	}
	node.TotalNodeNum = len(nodes)
	node.UnReadyNodeNum = unreadyNodes
	node.ReadyNodeNum = readyNodes
	node.CPUAllocatable = totalCpu
	node.CPURequested = usedCpu
	node.MemoryAllocatable = totalMemory
	node.MemoryRequested = usedMemory
	return &node
}
