package v1

import (
	"cmp/api/v1/response"
	"cmp/common"
	"cmp/model/k8s"
	"cmp/pkg"
	"cmp/pkg/cluster"
	"cmp/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateCluster(c *gin.Context) {
	d := k8s.Cluster{}
	err := c.ShouldBindJSON(&d)
	if err != nil {
		return
	}

	client, err := pkg.GetK8sClient(d.ConfigFileContentStr)
	if err != nil {
		response.FailWithMessage(response.CreateK8SClusterError, err.Error(), c)
		return
	}
	version, err := cluster.GetClusterStatus(client)
	if err != nil {
		response.FailWithMessage(response.CreateK8SClusterError, "连接集群异常,请检查网络", c)

	}
	d.ClusterStatus = true
	fmt.Println(d.ClusterStatus)
	d.ClusterVersion = version
	number, err := cluster.GetClusterNumber(client)
	if err != nil {
		response.FailWithMessage(response.CreateK8SClusterError, "获取集群节点数量异常", c)
	}
	d.NodeNumber = number
	if err := service.CreateCluster(d); err != nil {
		common.Log.Error("创建集群失败", zap.Any("err", err))
		response.FailWithMessage(response.CreateK8SClusterError, "创建集群失败", c)
		return
	} else {

		response.OkWithMessage("创建集群成功", c)
		return
	}

}

func ListCluster(c *gin.Context) {

	query := k8s.PaginationQ{}
	if c.ShouldBindQuery(&query) != nil {
		response.FailWithMessage(response.ParamError, response.ParamErrorMsg, c)
		return
	}
	var K8sCluster []k8s.Cluster
	if err := service.ListCluster(&query, &K8sCluster); err != nil {
		common.Log.Error("获取集群失败", zap.Any("err", err))
		response.FailWithMessage(response.InternalServerError, "获取集群失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			Data:  K8sCluster,
			Total: query.Total,
			Size:  query.Size,
			Page:  query.Page,
		}, "获取集群成功", c)
	}
}
func DelCluster(c *gin.Context) {

	var id k8s.ClusterIds
	err := c.ShouldBindJSON(&id)
	if err != nil {
		return
	}
	err2 := service.DelCluster(id)
	if err2 != nil {
		response.FailWithMessage(response.InternalServerError, "删除失败", c)
		return
	}
	response.Ok(c)
	return

}
func GetK8SClusterDetail(c *gin.Context) {

	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	data := cluster.GetClusterInfo(client)
	response.OkWithDetailed(data, "操作成功", c)

}
