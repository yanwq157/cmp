package pods

import (
	"cmp/model/k8s"
	k8scommon "cmp/pkg/common"
	v1 "k8s.io/api/core/v1"
)

type PodList struct {
	ListMet k8s.ListMeta             `json:"listMet"`
	Status  k8scommon.ResourceStatus `json:"status"`
	Pods    []Pod                    `json:"pods"`
}

type PodStatus struct {
	Status          string              `json:"status"`
	PodPhase        v1.PodPhase         `json:"podPhase"`
	ContainerStatus []v1.ContainerState `json:"containerStatus"`
}

type Pod struct {
	ObjectMeta      k8s.ObjectMeta    `json:"objectMeta"`
	TypeMeta        k8s.TypeMeta      `json:"typeMeta"`
	Status          string            `json:"status"`
	RestartCount    int32             `json:"restartCount"`
	Warnings        []k8scommon.Event `json:"warnings"`
	NodeName        string            `json:"nodeName"`
	ContainerImages []string          `json:"containerImages"`
	PodIP           string            `json:"podIP"`
}

func ToPod(pod *v1.Pod, warnings []k8scommon.Event) Pod {
	podDetail := Pod{
		ObjectMeta:      k8s.NewObjectMeta(pod.ObjectMeta),
		TypeMeta:        k8s.NewTypeMeta(k8s.ResourceKindPod),
		Warnings:        warnings,
		Status:          getPodStatus(*pod),
		RestartCount:    getRestartCount(*pod),
		NodeName:        pod.Spec.NodeName,
		ContainerImages: k8scommon.GetContainerImages(&pod.Spec),
		PodIP:           pod.Status.PodIP,
	}
	return podDetail
}
