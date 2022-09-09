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
