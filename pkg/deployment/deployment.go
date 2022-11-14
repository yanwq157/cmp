package deployment

import (
	"cmp/common"
	"cmp/model/k8s"
	k8scommon "cmp/pkg/common"
	"cmp/pkg/event"
	"cmp/tools"
	"context"
	"fmt"
	"go.uber.org/zap"
	apps "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	deploymentutil "k8s.io/kubectl/pkg/util/deployment"
	"time"
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
	// 更新node数量, filteredTotalcurl 过滤后的数量
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

func DeleteCollectionDeployment(client *kubernetes.Clientset, deploymentList []k8s.RemoveDeploymentData) (err error) {
	common.Log.Info("批量删除deployment开始")
	for _, v := range deploymentList {
		common.Log.Info(fmt.Sprintf("delete deployment：%v, ns: %v", v.DeploymentName, v.Namespace))
		err := client.AppsV1().Deployments(v.Namespace).Delete(
			context.TODO(),
			v.DeploymentName,
			metav1.DeleteOptions{},
		)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}
	common.Log.Info("删除deployment已完成")
	return nil
}

func DeleteDeployment(client *kubernetes.Clientset, ns string, deploymentName string) (err error) {
	common.Log.Info(fmt.Sprintf("请求删除单个deployment：%v,namespace:%v", deploymentName, ns))
	return client.AppsV1().Deployments(ns).Delete(
		context.TODO(),
		deploymentName,
		metav1.DeleteOptions{},
	)
}

func ScaleDeployment(client *kubernetes.Clientset, ns string, deploymentName string, scaleNumber int32) (err error) {
	common.Log.Info(fmt.Sprintf("start scale of %v deployment in %v namespace", deploymentName, ns))
	scaleData, err := client.AppsV1().Deployments(ns).GetScale(
		context.TODO(),
		deploymentName,
		metav1.GetOptions{},
	)
	common.Log.Info(fmt.Sprintf("The deployment has changed from %v to %v", scaleData.Spec.Replicas, scaleNumber))

	scale := autoscalingv1.Scale{
		TypeMeta:   scaleData.TypeMeta,
		ObjectMeta: scaleData.ObjectMeta,
		Spec:       autoscalingv1.ScaleSpec{Replicas: scaleNumber},
		Status:     scaleData.Status,
	}
	_, err = client.AppsV1().Deployments(ns).UpdateScale(
		context.TODO(),
		deploymentName,
		&scale,
		metav1.UpdateOptions{})
	if err != nil {
		common.Log.Error("扩缩容出现异常", zap.Any("err:", err))
		return err
	}
	return nil
}

func RestartDeployment(client *kubernetes.Clientset, deploymentName string, namespace string) (err error) {
	common.Log.Info(fmt.Sprintf("下发应用重启指令，名称空间：%v，无状态应用：%v", namespace, deploymentName))
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
	_, err = client.AppsV1().Deployments(namespace).Patch(
		context.Background(),
		deploymentName,
		types.StrategicMergePatchType,
		[]byte(data),
		metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		common.Log.Error("应用重启失败", zap.Any("err:", err))
		return err
	}
	return nil
}

func RollDeployment(client *kubernetes.Clientset, deploymentName string, namespace string, reVersion int64) (err error) {
	common.Log.Info(fmt.Sprintf("应用：%v, 所属空间：%v, 版本回滚到%v", deploymentName, namespace, reVersion))
	if reVersion < 0 {
		return revisionNotFoundErr(reVersion)
	}
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to retrieve Deployment %s:%v", deploymentName, err)
	}
	if deployment.Spec.Paused {
		return fmt.Errorf("skipped rollback (deployment \"%s\" is paused)", deployment.Name)
	}
	if reVersion == 0 {
		common.Log.Warn("传递回滚版本号是：0，默认回退上一次版本")
		rsForRevision, err := deploymentRevision(deployment, client, reVersion)
		if err != nil {
			return err
		}

		for k := range rsForRevision.Annotations {
			if k == "deployment.kubernetes.io/revision" {
				deployment.Spec.Template = rsForRevision.Spec.Template
				if _, rollbackErr := client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{}); rollbackErr != nil {
					common.Log.Error("版本回退失败", zap.Any("err:", err))
					return rollbackErr
				}
				common.Log.Info("The rollback task was executed successfully")
				return nil
			}
		}
	}

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	replicaSetList, err := client.AppsV1().ReplicaSets(namespace).List(context.TODO(), options)
	if err != nil {
		return err
	}
	if len(replicaSetList.Items) <= 1 {
		return revisionNotFoundErr(reVersion)
	}

	for _, v := range replicaSetList.Items {
		currentVersion := tools.ParseStringToInt64(v.Annotations["deployment.kubernetes.io/revision"])
		common.Log.Info(fmt.Sprintf("currentVersion: %v", currentVersion))
		common.Log.Info(fmt.Sprintf("reVersion: %v", reVersion))

		if currentVersion == reVersion {
			deployment.Spec.Template = v.Spec.Template
			if _, rollbackErr := client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{}); rollbackErr != nil {
				common.Log.Error("版本回退失败", zap.Any("err:", err))
				return rollbackErr
			}
			common.Log.Info("The rollback task was executed successfully")
			return nil
		}
	}
	return nil
}

func deploymentRevision(deployment *apps.Deployment, c kubernetes.Interface, toRevision int64) (revision *apps.ReplicaSet, err error) {
	_, allOldRSs, newRS, err := deploymentutil.GetAllReplicaSets(deployment, c.AppsV1())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve replica sets from deployment %s:%v", deployment.Name, err)
	}
	allRSs := allOldRSs
	if newRS != nil {
		allRSs = append(allRSs, newRS)
	}
	var (
		latestReplicaSet   *apps.ReplicaSet
		latestRevision     = int64(-1)
		previousReplicaSet *apps.ReplicaSet
		previousRevision   = int64(-1)
	)
	for _, rs := range allRSs {
		if v, err := deploymentutil.Revision(rs); err == nil {
			common.Log.Info(fmt.Sprintf("v: %v", v))

			if toRevision == 0 {
				if latestRevision < v {
					previousRevision = latestRevision
					previousReplicaSet = latestReplicaSet
					latestRevision = v
					latestReplicaSet = rs
				} else if previousRevision < v {
					previousRevision = v
					previousReplicaSet = rs
				}
			} else if toRevision == v {
				return rs, nil
			}
		}
	}
	if toRevision > 0 {
		return nil, revisionNotFoundErr(toRevision)
	}
	if previousReplicaSet == nil {
		return nil, fmt.Errorf("no rollout history found for deployment %q", deployment.Name)
	}

	return previousReplicaSet, nil
}

func revisionNotFoundErr(r int64) error {
	common.Log.Warn("没有找到可回滚版本")
	return fmt.Errorf("unable to find specified revision %v in history", r)
}
