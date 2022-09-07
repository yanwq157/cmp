package v1

import (
	"cmp/api/v1/response"
	"cmp/common"
	cluster2 "cmp/model/cluster"
	"cmp/pkg/cluster"
	"cmp/pkg/init"
	"cmp/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateCluster(c *gin.Context) {
	d := cluster2.Cluster{}
	err := c.ShouldBindJSON(&d)
	if err != nil {
		return
	}
	fmt.Println(d.ConfigFileContentStr)

	client, err := init.GetK8sClient(d.ConfigFileContentStr)
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

	query := cluster2.PaginationQ{}
	if c.ShouldBindQuery(&query) != nil {
		response.FailWithMessage(response.ParamError, response.ParamErrorMsg, c)
		return
	}
	var K8sCluster []cluster2.Cluster
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

	var id cluster2.ClusterIds
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

	client, err := init.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	data := cluster.GetClusterInfo(client)
	response.OkWithDetailed(data, "操作成功", c)

}
