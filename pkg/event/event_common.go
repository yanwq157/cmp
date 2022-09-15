package event

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
