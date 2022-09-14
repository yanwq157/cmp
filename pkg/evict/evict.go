package evict

import (
	"cmp/common"
	"context"
	"fmt"
	policy "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	systemNamespace = "kube-system"
)

//驱逐节点上不是kube-system命名空间中的所有pod
func EvictsNodePods(clinet *kubernetes.Clientset, nodeName string) error {
	pods, err := clinet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})

	if err != nil {
		return err
	}
	for _, i := range pods.Items {
		if i.Namespace == systemNamespace {
			continue
		} else {
			common.Log.Info(fmt.Sprintf("开始驱逐Node：%v,节点Namespace：%v下的pod：%v", nodeName, i.Namespace, i.Name))
			err := EvictsPod(clinet, i.Name, i.Namespace)
			if err != nil {
				common.Log.Error(fmt.Sprintf("驱逐Pod%v失败", i.Name))
			}
		}
	}
	common.Log.Info(fmt.Sprintf("已成功从节点：%v中驱逐所有pod", nodeName))
	return nil
}

// EvictsPod 驱逐pod
func EvictsPod(client *kubernetes.Clientset, name, namespace string) error {
	var gracePeriodSeconds int64 = 0
	propagationPolicy := metav1.DeletePropagationForeground
	//立即删除，并删除依赖
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &propagationPolicy,
	}
	return client.PolicyV1beta1().Evictions(namespace).Evict(context.TODO(), &policy.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		DeleteOptions: deleteOptions,
	})
}
