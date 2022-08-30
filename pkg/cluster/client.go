package pkg

import (
	"cmp/common"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetK8sClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "configs/config")
	if err != nil {
		common.Log.Error("KubeConfig内容错误", zap.Any("err", err))
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.Log.Error("创建Client失败", zap.Any("err", err))
		return nil, err
	}
	return clientSet, err
}
