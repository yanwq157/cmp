package deployment

import (
	"cmp/common"
	k8scommon "cmp/pkg/common"
	"cmp/pkg/event"
	"cmp/pkg/service"
	"cmp/tools"
	"context"
	"fmt"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"sort"
)

type RollingUpdateStrategy struct {
	MaxSurgr       *intstr.IntOrString `json:"maxSurgr"`
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable"`
}

type StatusInfo struct {
	Replicas    int32 `json:"replicas"`
	Updated     int32 `json:"updated"`
	Avilable    int32 `json:"avilable"`
	Unavailable int32 `json:"unavailable"`
}

type DeploymentDetail struct {
	Deployment            `json:"deployment"`
	Selector              map[string]string `json:"selector"`
	StatusInfo            `json:"statusInfo"`
	Conditions            []k8scommon.Condition       `json:"conditions"`
	Strategy              apps.DeploymentStrategyType `json:"strategy"`
	MinReadySeconds       int32                       `json:"minReadySeconds"`
	RollingUpdateStrategy *RollingUpdateStrategy      `json:"rollingUpdateStrategy"`
	RevisionHistoryLimit  *int32                      `json:"revisionHistoryLimit"`
	Events                []v1.Event                  `json:"events"`
	HistoryVersion        []HistoryVersion            `json:"historyVersion"`
	PodList               *PodList                    `json:"podList"`
	SvcList               *service.ServiceList        `json:"svcList"`
}

type HistoryVersion struct {
	CreateTime metaV1.Time `json:"createTime"`
	Image      string      `json:"image"`
	Version    int64       `json:"version"`
	Namespace  string      `json:"namespace"`
	Name       string      `json:"name"`
}

func GetDeploymentDetail(client *kubernetes.Clientset, namespace string, deploymentName string) (*DeploymentDetail, error) {
	common.Log.Info(fmt.Sprintf("Getting details of %v deployment in %v namespace", deploymentName, namespace))

	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	selector, err := metaV1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metaV1.ListOptions{LabelSelector: selector.String()}

	channels := &k8scommon.ResourceChannels{
		ReplicaSetList: k8scommon.GetReplicaSetListChannelWithOptions(client, k8scommon.NewSameNamespaceQuery(namespace), options, 1),
		PodList:        k8scommon.GetPodListChannelWithOptions(client, k8scommon.NewSameNamespaceQuery(namespace), options, 1),
		EventList:      k8scommon.GetEventListChannelWithOptions(client, k8scommon.NewSameNamespaceQuery(namespace), options, 1),
	}
	rawRs := <-channels.ReplicaSetList.List
	err = <-channels.ReplicaSetList.Error
	if err != nil {
		return nil, err
	}
	rawPods := <-channels.PodList.List
	err = <-channels.PodList.Error
	if err != nil {
		return nil, err
	}
	rawEvents := <-channels.EventList.List
	err = <-channels.EventList.Error
	if err != nil {
		return nil, err
	}

	var rollingUpdateStrategy *RollingUpdateStrategy
	if deployment.Spec.Strategy.RollingUpdate != nil {
		rollingUpdateStrategy = &RollingUpdateStrategy{
			MaxSurgr:       deployment.Spec.Strategy.RollingUpdate.MaxSurge,
			MaxUnavailable: deployment.Spec.Strategy.RollingUpdate.MaxUnavailable,
		}
	}
	events, _ := event.GetEvents(client, namespace, fmt.Sprintf("involvedObject.name=%v", deploymentName))
	serviceList, _ := service.GetToService(client, namespace, deploymentName)
	return &DeploymentDetail{
		Deployment:            toDeployment(deployment, rawRs.Items, rawPods.Items, rawEvents.Items),
		Selector:              deployment.Spec.Selector.MatchLabels,
		StatusInfo:            GetStatusInfo(&deployment.Status),
		Conditions:            getConditions(deployment.Status.Conditions),
		Strategy:              deployment.Spec.Strategy.Type,
		MinReadySeconds:       deployment.Spec.MinReadySeconds,
		RollingUpdateStrategy: rollingUpdateStrategy,
		RevisionHistoryLimit:  deployment.Spec.RevisionHistoryLimit,
		Events:                events,
		PodList:               getDeploymentToPod(client, deployment),
		SvcList:               serviceList,
		HistoryVersion:        getDeploymentHistory(namespace, deploymentName, rawRs.Items),
	}, nil
}

func GetStatusInfo(deploymentStatus *apps.DeploymentStatus) StatusInfo {
	return StatusInfo{
		Replicas:    deploymentStatus.Replicas,
		Updated:     deploymentStatus.UpdatedReplicas,
		Avilable:    deploymentStatus.AvailableReplicas,
		Unavailable: deploymentStatus.UnavailableReplicas,
	}
}

func getDeploymentHistory(namespace string, deploymentName string, rs []apps.ReplicaSet) []HistoryVersion {
	var historyVersion []HistoryVersion
	for _, v := range rs {
		if namespace == v.Namespace && deploymentName == v.OwnerReferences[0].Name {
			history := HistoryVersion{
				CreateTime: v.CreationTimestamp,
				Image:      v.Spec.Template.Spec.Containers[0].Image,
				Version:    tools.ParseStringToInt64(v.Annotations["deployment.kubernetes.io/revision"]),
				Namespace:  v.Namespace,
				Name:       v.OwnerReferences[0].Name,
			}
			historyVersion = append(historyVersion, history)
		}
	}
	sort.Sort(historiesByRevision(historyVersion))
	return historyVersion
}

type historiesByRevision []HistoryVersion

func (h historiesByRevision) Len() int {
	return len(h)
}
func (h historiesByRevision) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h historiesByRevision) Less(i, j int) bool {
	return h[j].Version < h[i].Version
}
