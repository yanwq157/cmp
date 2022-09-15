package v1

import (
	"cmp/api/v1/response"
	"cmp/pkg"
	"cmp/pkg/deployment"
	"cmp/pkg/parser"
	"github.com/gin-gonic/gin"
)

func GetDeploymentList(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	//后续分页
	//解析路径参数中的命名空间
	//不传or传不存在的输出所有空间
	nsQuery := parser.ParseNamespacePathParameter(c)
	data, err := deployment.GetDeploymentList(client, nsQuery)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "操作成功", c)
	return

}
