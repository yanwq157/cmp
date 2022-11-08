package deployment

import (
	k8scommon "cmp/pkg/common"
	"cmp/pkg/event"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getStatus(list *apps.DeploymentList, rs []apps.ReplicaSet, pods []v1.Pod, events []v1.Event) k8scommon.ResourceStatus {
	info := k8scommon.ResourceStatus{}
	if list == nil {
		return info
	}
	for _, deployment := range list.Items {
		matchingPods := k8scommon.FilterDeploymentPodsByOwnerReference(deployment, rs, pods)
		podInfo := k8scommon.GetPodInfo(deployment.Status.Replicas, deployment.Spec.Replicas, matchingPods)
		warnings := event.GetPodsEventWarnings(events, matchingPods)

		if len(warnings) > 0 {
			info.Failed++
		} else if podInfo.Pending > 0 {
			info.Pending++
		} else {
			info.Running++
		}
	}
	return info
}

func getConditions(deploymentConditions []apps.DeploymentCondition) []k8scommon.Condition {
	conditions := make([]k8scommon.Condition, 0)
	for _, condition := range deploymentConditions {
		conditions = append(conditions, k8scommon.Condition{
			Type:               string(condition.Type),
			Status:             metaV1.ConditionStatus(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime,
			LastProbeTime:      condition.LastUpdateTime,
		})
	}
	return conditions
}
