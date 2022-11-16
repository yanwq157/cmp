package v1

import (
	"cmp/api/v1/response"
	"cmp/pkg"
	"cmp/pkg/parser"
	"cmp/pkg/pods"
	"github.com/gin-gonic/gin"
)

func GetPodsListController(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	nsQuery := parser.ParseNamespacePathParameter(c)

	data, err := pods.GetPodList(client, nsQuery)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithData(data, c)
	return
}
