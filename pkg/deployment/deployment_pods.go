package deployment

import (
	"cmp/common"
	"cmp/model/k8s"
	k8scommon "cmp/pkg/common"
	"cmp/pkg/event"
	"cmp/pkg/pods"
	"context"
	"go.uber.org/zap"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodList struct {
	ListMeta k8s.ListMeta             `json:"listMeta"`
	Status   k8scommon.ResourceStatus `json:"status"`
	Pods     []pods.Pod               `json:"pods"`
}

func getDeploymentToPod(client *kubernetes.Clientset, deployment *apps.Deployment) (po *PodList) {
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	podData, err := client.CoreV1().Pods(deployment.Namespace).List(context.TODO(), options)
	if err != nil {
		common.Log.Error("Get a pod execption from the deployment", zap.Any("err", err))
	}
	podList := PodList{
		Pods: make([]pods.Pod, 0),
	}
	podList.ListMeta = k8s.ListMeta{TotalItems: len(podData.Items)}
	for _, pod := range podData.Items {
		warnings := event.GetPodsEventWarnings(nil, []v1.Pod{})
		podDetail := pods.ToPod(&pod, warnings)
		podList.Pods = append(podList.Pods, podDetail)
	}
	return &podList
}
