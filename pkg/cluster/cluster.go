package pkg

import (
	"cmp/common"
	"context"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetClusterStatus(c *kubernetes.Clientset) (string, error) {
	version, err := c.ServerVersion()
	if err != nil {
		common.Log.Error("get version from cluster failed", zap.Any("err: ", err))
		return "", err
	}
	return version.String(), nil
}

func GetClusterNumber(c *kubernetes.Clientset) (int, error) {
	number, err := c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(number.Items), nil
}
