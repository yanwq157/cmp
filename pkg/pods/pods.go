package pods

import (
	"cmp/common"
	"cmp/model/k8s"
	k8scommon "cmp/pkg/common"
	"cmp/pkg/event"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

func GetPodList(client *kubernetes.Clientset, nsQuery *k8scommon.NamespaceQuery) (*PodList, error) {
	common.Log.Info("Getting list of all pods in the cluster")
	channels := &k8scommon.ResourceChannels{
		PodList:   k8scommon.GetPodListChannelWithOptions(client, nsQuery, metav1.ListOptions{}, 1),
		EventList: k8scommon.GetEventListChannel(client, nsQuery, 1),
	}
	return GetPodListFromChannels(channels)
}

func GetPodListFromChannels(channels *k8scommon.ResourceChannels) (*PodList, error) {
	pods := <-channels.PodList.List
	err := <-channels.PodList.Error
	if err != nil {
		return nil, err
	}

	eventList := <-channels.EventList.List
	err = <-channels.EventList.Error
	if err != nil {
		return nil, err
	}

	podList := ToPodList(pods.Items, eventList.Items)
	podList.Status = getStatus(pods, eventList.Items)
	return &podList, nil
}

func ToPodList(pods []v1.Pod, events []v1.Event) PodList {
	podList := PodList{
		Pods: make([]Pod, 0),
	}
	//podCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(pods), dsQuery)
	//pods = fromCells(podCells)
	//podList.ListMeta = k8s.ListMeta{TotalItems: filteredTotal}

	for _, pod := range pods {
		warnings := event.GetPodsEventWarnings(events, []v1.Pod{pod})
		podDetail := ToPod(&pod, warnings)
		podList.Pods = append(podList.Pods, podDetail)
	}
	return podList
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
