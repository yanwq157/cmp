package event

import (
	"cmp/pkg/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

func GetPodsEventWarnings(events []v1.Event, pods []v1.Pod) []common.Event {
	result := make([]common.Event, 0)
	events = GetWarningEvents(events)
	failedPods := make([]v1.Pod, 0)

	for _, pod := range pods {
		if !isReadyOrSucceeded(pod) {
			failedPods = append(failedPods, pod)
		}
	}
	events = filterEventsByPodsUID(events, failedPods)
	events = removeDuplicates(events)

	for _, event := range events {
		result = append(result, common.Event{
			Message: event.Message,
			Reason:  event.Reason,
			Type:    event.Type,
		})
	}
	return result
}

func GetWarningEvents(events []v1.Event) []v1.Event {
	return filterEventsByType(FillEventsType(events), v1.EventTypeWarning)
}

func isReadyOrSucceeded(pod v1.Pod) bool {
	if pod.Status.Phase == v1.PodSucceeded {
		return true
	}
	if pod.Status.Phase == v1.PodRunning {
		for _, c := range pod.Status.Conditions {
			if c.Type == v1.PodReady {
				if c.Status == v1.ConditionFalse {
					return false
				}
			}
			return true
		}
	}
	return false
}

func filterEventsByPodsUID(events []v1.Event, pods []v1.Pod) []v1.Event {
	result := make([]v1.Event, 0)
	podEventMap := make(map[types.UID]bool, 0)
	if len(pods) == 0 || len(events) == 0 {
		return result
	}
	for _, pod := range pods {
		podEventMap[pod.UID] = true
	}
	for _, event := range events {
		if _, exists := podEventMap[event.InvolvedObject.UID]; exists {
			result = append(result, event)
		}
	}
	return result
}

func removeDuplicates(slice []v1.Event) []v1.Event {
	visited := make(map[string]bool, 0)
	result := make([]v1.Event, 0)
	for _, elem := range slice {
		if !visited[elem.Reason] {
			visited[elem.Reason] = true
			result = append(result, elem)
		}
	}
	return result
}

func FillEventsType(events []v1.Event) []v1.Event {
	for i := range events {
		if len(events[i].Type) == 0 {
			if isFailedReason(events[i].Reason, FailedReasonPartials...) {
				events[i].Type = v1.EventTypeWarning
			} else {
				events[i].Type = v1.EventTypeNormal
			}
		}
	}
	return events
}

func filterEventsByType(events []v1.Event, eventType string) []v1.Event {
	if len(eventType) == 0 || len(events) == 0 {
		return events
	}
	result := make([]v1.Event, 0)
	for _, event := range events {
		if event.Type == eventType {
			result = append(result, event)
		}
	}
	return result
}

func isFailedReason(reason string, partials ...string) bool {
	for _, partial := range partials {
		if strings.Contains(strings.ToLower(reason), partial) {
			return true
		}
	}
	return false
}
