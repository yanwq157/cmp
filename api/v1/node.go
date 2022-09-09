package v1

import (
	"cmp/api/v1/response"
	"cmp/pkg"
	"cmp/pkg/node"
	"github.com/gin-gonic/gin"
)

func GetNodes(c *gin.Context) {

	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	data, err := node.GetNodeList(client)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "操作成功", c)
	return
}

func GetNodeDetail(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	name := c.Query("name")
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
	}
	data, err := node.GetNodeDetail(client, name)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "操作成功", c)
	return
}

type Status struct {
	NodeName    string `json:"node_name"`
	Unscheduled bool   `json:"unscheduled"`
}

func NodeUnschedulable(c *gin.Context) {
	var nodeUnscheduled Status
	err := c.ShouldBindJSON(&nodeUnscheduled)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	data, err := node.NodeUnschedulable(client, nodeUnscheduled.NodeName, nodeUnscheduled.Unscheduled)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "操作成功", c)
	return
}

func CordonNode(c *gin.Context) {
	nodeName := c.Query("node_name")
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	if ok, err := node.CordonNode(client, nodeName); !ok {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.Ok(c)
}
