package event

import (
	k8scommon "cmp/pkg/common"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func GetNodeEvents(client *kubernetes.Clientset, nodeName string) (*v1.EventList, error) {
	events, err := client.CoreV1().Events(v1.NamespaceAll).List(context.TODO(),
		metaV1.ListOptions{FieldSelector: fmt.Sprintf("involvedObject.name=%v", nodeName)})

	if err != nil {
		return nil, err
	}
	return events, nil
}

var FailedReasonPartials = []string{"failed", "err", "exceeded", "invalid", "unhealthy",
	"mismatch", "insufficient", "conflict", "outof", "nil", "backoff"}

func GetEvents(client *kubernetes.Clientset, namespace, resourceName string) ([]v1.Event, error) {
	fieldSelector, err := fields.ParseSelector("involvedObject.name" + "=" + resourceName)
	if err != nil {
		return nil, err
	}

	channles := &k8scommon.ResourceChannels{
		EventList: k8scommon.GetEventListChannelWithOptions(client, k8scommon.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labels.Everything().String(),
				FieldSelector: fieldSelector.String(),
			},
			1),
	}
	eventList := <-channles.EventList.List
	if err := <-channles.EventList.Error; err != nil {
		return nil, err

	}
	return FillEventsType(eventList.Items), nil
}
