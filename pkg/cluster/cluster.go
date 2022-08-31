package pkg

import (
	"cmp/common"
	"go.uber.org/zap"
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

func GetClusterNumber(c *kubernetes.Clientset) (string, error) {
	number, err := c.CoreV1().Nodes().Get()
}
