package v1

import (
	"cmp/api/v1/response"
	"cmp/pkg"
	"cmp/pkg/namespace"
	"github.com/gin-gonic/gin"
)

func GetNamespaceList(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	namespace, err := namespace.GetNamespaceList(client)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithData(namespace, c)
	return
}
