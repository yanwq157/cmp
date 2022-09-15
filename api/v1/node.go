package v1

import (
	"cmp/api/v1/response"
	"cmp/pkg"
	"cmp/pkg/node"
	"github.com/gin-gonic/gin"
	"net/http"
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
	name := c.Query("nodeName")
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
	NodeName    string `json:"nodeName"`
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
	data, err := node.K8sNodeUnschedulable(client, nodeUnscheduled.NodeName, nodeUnscheduled.Unscheduled)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "操作成功", c)
	return
}

func CordonNode(c *gin.Context) {
	nodeName := c.Query("nodeName")
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

func RemoveNode(c *gin.Context) {
	nodeName := c.Query("nodeName")
	if nodeName == "" {
		response.FailWithMessage(http.StatusNotFound, "移除节点名称不能为空", c)
		return
	}

	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	go node.RemoveNode(client, nodeName)
	response.Ok(c)
}

type collectionNode struct {
	NodeName []string `json:"nodeName"`
}

func CollectionNodeUnschedule(c *gin.Context) {
	var nodeUnscheduled collectionNode
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
	err = node.CollectionNodeUnschedule(client, nodeUnscheduled.NodeName)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

func CollectionCordonNode(c *gin.Context) {
	var nodeUnscheduled collectionNode
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
	err2 := node.CollectionCordonNode(client, nodeUnscheduled.NodeName)
	if err2 != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.Ok(c)
}
