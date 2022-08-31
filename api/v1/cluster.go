package v1

import (
	"cmp/api/v1/response"
	"cmp/common"
	"cmp/model"
	pkg "cmp/pkg/cluster"
	"cmp/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddCluster(c *gin.Context) {
	d := model.Cluster{}
	err := c.ShouldBindJSON(&d)
	if err != nil {
		return
	}
	client, err := pkg.GetK8sClient()

	version, err := pkg.GetClusterStatus(client)
	if err != nil {
		response.FailWithMessage(response.CreateK8SClusterError, "连接集群异常,请检查网络", c)

	}
	d.ClusterVersion = version
	number, err := pkg.GetClusterNumber(client)
	if err != nil {
		response.FailWithMessage(response.CreateK8SClusterError, "获取集群节点数量异常", c)
	}
	if err := service.CreateCluster(d); err != nil {
		common.Log.Error("创建集群失败", zap.Any("err", err))
		response.FailWithMessage(response.CreateK8SClusterError, "创建集群失败", c)
		return
	} else {
		response.OkWithMessage("创建集群成功", c)
		return
	}
}
