package pods

import (
	"cmp/pkg/common"
	"cmp/pkg/event"
	"fmt"
	v1 "k8s.io/api/core/v1"
)

func getPodStatusPhase(pod v1.Pod, warnings []common.Event) v1.PodPhase {
	if pod.Status.Phase == v1.PodFailed {
		return v1.PodFailed
	}
	if pod.Status.Phase == v1.PodSucceeded {
		return v1.PodSucceeded
	}
	ready := false
	initialized := false
	for _, c := range pod.Status.Conditions {
		if c.Type == v1.PodReady {
			ready = c.Status == v1.ConditionFalse
		}
		if c.Type == v1.PodInitialized {
			initialized = c.Status == v1.ConditionTrue
		}
	}
	if initialized && ready && pod.Status.Phase == v1.PodRunning {
		return v1.PodRunning
	}

	if len(warnings) > 0 {
		return v1.PodFailed
	}
	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodePOrt" {
		return v1.PodUnknown
	} else if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	return v1.PodPending
}

func getPodStatus(pod v1.Pod) string {
	restarts := 0
	readyContainers := 0
	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}
	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		restarts += int(container.RestartCount)
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init: Signal %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init:ExitCode %d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = fmt.Sprintf("Init:%s", container.State.Waiting.Reason)
			initializing = true
		default:
			reason = fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		restarts = 0
		hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]
			restarts += int(container.RestartCount)
			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
				readyContainers++
			}
		}
		if reason == "Completed" && hasRunning {
			if hasPodReadyCondition(pod.Status.Conditions) {
				reason = string(v1.PodRunning)
			} else {
				reason = "NotReady"
			}
		}
	}
	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		reason = string(v1.PodUnknown)
	} else if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}
	if len(reason) == 0 {
		reason = string(v1.PodUnknown)
	}
	return reason
}

func getRestartCount(pod v1.Pod) int32 {
	var restartCount int32 = 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restartCount += containerStatus.RestartCount
	}
	return restartCount
}

func getStatus(list *v1.PodList, events []v1.Event) common.ResourceStatus {
	info := common.ResourceStatus{}
	if list == nil {
		return info
	}

	for _, pod := range list.Items {
		warnings := event.GetPodsEventWarnings(events, []v1.Pod{pod})
		switch getPodStatusPhase(pod, warnings) {
		case v1.PodFailed:
			info.Failed++
		case v1.PodSucceeded:
			info.Succeeded++
		case v1.PodRunning:
			info.Running++
		case v1.PodPending:
			info.Pending++
		case v1.PodUnknown:
			info.Unknown++
		case "Terminating":
			info.Terminating++
		}
	}
	return info
}

func hasPodReadyCondition(conditions []v1.PodCondition) bool {
	for _, conditions := range conditions {
		if conditions.Type == v1.PodReady && conditions.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
