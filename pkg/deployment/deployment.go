package deployment

import (
	"cmp/common"
	"cmp/model/k8s"
	k8scommon "cmp/pkg/common"
	"cmp/pkg/event"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentList struct {
	ListMeta    k8s.ListMeta             `json:"listMeta"`
	Status      k8scommon.ResourceStatus `json:"status"`
	Deployments []Deployment             `json:"deployments"`
}

type Deployment struct {
	ObjectMeta          k8s.ObjectMeta    `json:"objectMeta"`
	TypeMeta            k8s.TypeMeta      `json:"typeMeta"`
	Pods                k8scommon.PodInfo `json:"pods"`
	ComtainerImages     []string          `json:"comtainerImages"`
	InitContainerImages []string          `json:"initContainerImages"`
	DeploymentStatus    DeploymentStatus  `json:"deploymentStatus"`
}

type DeploymentStatus struct {
	Replicas            int32 `json:"replicas"`
	UpdatedReplicas     int32 `json:"updatedReplicas"`
	ReadyReplicas       int32 `json:"readyReplicas"`
	AvailableReplicas   int32 `json:"availableReplicas"`
	UnavailableReplicas int32 `json:"unavailableReplicas"`
}

func GetDeploymentList(client *kubernetes.Clientset, nsQuery *k8scommon.NamespaceQuery) (*DeploymentList, error) {
	common.Log.Info("获取集群中deployment")
	channels := &k8scommon.ResourceChannels{
		DeploymentList: k8scommon.GetDeploymentListChannel(client, nsQuery, 1),
		PodList:        k8scommon.GetPodListChannel(client, nsQuery, 1),
		EventList:      k8scommon.GetEventListChannel(client, nsQuery, 1),
		ReplicaSetList: k8scommon.GetReplicaSetListChannel(client, nsQuery, 1),
	}
	return GetDeploymentListFromChannels(channels)
}

func GetDeploymentListFromChannels(channels *k8scommon.ResourceChannels) (*DeploymentList, error) {
	deployments := <-channels.DeploymentList.List
	err := <-channels.DeploymentList.Error
	if err != nil {
		return nil, err
	}
	pods := <-channels.PodList.List
	err = <-channels.PodList.Error
	if err != nil {
		return nil, err
	}
	events := <-channels.EventList.List
	err = <-channels.EventList.Error
	if err != nil {
		return nil, err
	}
	rs := <-channels.ReplicaSetList.List
	err = <-channels.ReplicaSetList.Error
	if err != nil {
		return nil, err
	}
	deploymentList := toDeploymentList(deployments.Items, pods.Items, events.Items, rs.Items)
	deploymentList.Status = getStatus(deployments, rs.Items, pods.Items, events.Items)
	return deploymentList, nil
}

func toDeploymentList(deployments []apps.Deployment, pods []v1.Pod, events []v1.Event, rs []apps.ReplicaSet) *DeploymentList {
	deploymentList := &DeploymentList{
		Deployments: make([]Deployment, 0),
		ListMeta:    k8s.ListMeta{TotalItems: len(deployments)},
	}
	// 解析前端传递的参数, filterBy=name,1.1&itemsPerPage=10&name=&namespace=default&page=1&sortBy=d,creationTimestamp
	// sortBy=d 倒序, sortBy=a 正序, 排序按照a-z
	//dataSelect := parser.ParseDataSelectPathParameter(dsQuery)
	// 过滤
	//nodeCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(nodes), dataSelect)
	//nodes = fromCells(nodeCells)
	// 更新node数量, filteredTotal过滤后的数量
	//nodeList.ListMeta = k8s.ListMeta{TotalItems: filteredTotal}
	//deploymentCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(deployments), dsQuery)
	//deployments = fromCells(deploymentCells)
	//deploymentList.ListMeta = k8s.ListMeta{TotalItems: filteredTotal}
	for _, deployment := range deployments {
		deploymentList.Deployments = append(deploymentList.Deployments, toDeployment(&deployment, rs, pods, events))
	}
	return deploymentList
}

func toDeployment(deployment *apps.Deployment, rs []apps.ReplicaSet, pod []v1.Pod, events []v1.Event) Deployment {
	matchingPods := k8scommon.FilterDeploymentPodsByOwnerReference(*deployment, rs, pod)
	podInfo := k8scommon.GetPodInfo(deployment.Status.Replicas, deployment.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(events, matchingPods)
	return Deployment{
		ObjectMeta:          k8s.NewObjectMeta(deployment.ObjectMeta),
		TypeMeta:            k8s.NewTypeMeta(k8s.ResourceKindDeployment),
		Pods:                podInfo,
		ComtainerImages:     k8scommon.GetContainerImages(&deployment.Spec.Template.Spec),
		InitContainerImages: k8scommon.GetInitContainerImages(&deployment.Spec.Template.Spec),
		DeploymentStatus:    getDeploymentStatus(deployment),
	}
}
func getDeploymentStatus(deployment *apps.Deployment) DeploymentStatus {
	return DeploymentStatus{
		Replicas:            deployment.Status.Replicas,
		UpdatedReplicas:     deployment.Status.UpdatedReplicas,
		ReadyReplicas:       deployment.Status.ReadyReplicas,
		AvailableReplicas:   deployment.Status.AvailableReplicas,
		UnavailableReplicas: deployment.Status.UnavailableReplicas,
	}
}
