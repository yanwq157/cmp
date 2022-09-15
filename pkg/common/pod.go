package common

import (
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FilterDeploymentPodsByOwnerReference(deployment apps.Deployment, allRS []apps.ReplicaSet, allPods []v1.Pod) []v1.Pod {
	var matchingsPods []v1.Pod
	for _, rs := range allRS {
		if metav1.IsControlledBy(&rs, &deployment) {
			//无法将 'FilterPodsByControllerRef(&rs,allPods)' (类型 []v1.Pod) 用作类型 v1.Pod
			matchingsPods = append(matchingsPods, FilterPodsByControllerRef(&rs, allPods)...)
		}
	}
	return matchingsPods
}
func FilterPodsByControllerRef(owner metav1.Object, allPods []v1.Pod) []v1.Pod {
	var matchingPods []v1.Pod
	for _, pod := range allPods {
		if metav1.IsControlledBy(&pod, owner) {
			matchingPods = append(matchingPods, pod)
		}
	}
	return matchingPods
}

func GetContainerImages(podTemplate *v1.PodSpec) []string {
	var ContainerImages []string
	for _, Container := range podTemplate.Containers {
		ContainerImages = append(ContainerImages, Container.Image)
	}
	return ContainerImages
}

func GetInitContainerImages(podTemplate *v1.PodSpec) []string {
	var initContainerImages []string
	for _, initContainer := range podTemplate.InitContainers {
		initContainerImages = append(initContainerImages, initContainer.Image)
	}
	return initContainerImages
}
