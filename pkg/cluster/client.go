package cluster

import (
	"cmp/common"
	"cmp/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
)

func GetK8sClient(k8sConf string) (*kubernetes.Clientset, error) {
	//config, err := clientcmd.BuildConfigFromFlags("", "configs/config")
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(k8sConf))
	if err != nil {
		common.Log.Error("KubeConfig内容错误", zap.Any("err", err))
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.Log.Error("KubeConfig内容错误", zap.Any("err", err))
		return nil, err
	}
	return clientSet, err
}

// ClusterID 公共方法, 获取指定k8s集群的KubeConfig
func ClusterID(c *gin.Context) (*kubernetes.Clientset, error) {

	clusterId := c.DefaultQuery("clusterId", "1")
	clusterIdUint, err := strconv.ParseUint(clusterId, 10, 32)
	cluster, err := service.GetK8sClusterID(uint(clusterIdUint))
	if err != nil {
		common.Log.Error("获取集群失败", zap.Any("err", err))
		return nil, err
	}

	client, _ := GetK8sClient(cluster.ConfigFileContentStr)
	return client, nil
}
