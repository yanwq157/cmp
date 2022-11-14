package v1

import (
	"cmp/api/v1/response"
	"cmp/common"
	"cmp/model/k8s"
	"cmp/pkg"
	"cmp/pkg/deployment"
	"cmp/pkg/parser"
	"cmp/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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

func DeleteCollectionDeployment(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	var deploymentList []k8s.RemoveDeploymentData
	err = c.ShouldBindJSON(&deploymentList)
	if err != nil {
		response.FailWithMessage(http.StatusNotFound, err.Error(), c)
		return
	}
	err = deployment.DeleteCollectionDeployment(client, deploymentList)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

func DeleteDeployment(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	var deploymentData k8s.RemoveDeploymentToServiceData
	err2 := c.ShouldBindJSON(&deploymentData)
	if err2 != nil {
		response.FailWithMessage(http.StatusNotFound, err2.Error(), c)
		return
	}

	err = deployment.DeleteDeployment(client, deploymentData.Namespace, deploymentData.DeploymentName)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	common.Log.Info(fmt.Sprintf("deployment:%v已删除", deploymentData.DeploymentName))

	if deploymentData.IsDeleteService {
		serviceErr := service.DeleteService(client, deploymentData.Namespace, deploymentData.ServiceName)
		if serviceErr != nil {
			common.Log.Error("删除相关Service出错", zap.Any("err: ", serviceErr))
			response.FailWithMessage(response.InternalServerError, err.Error(), c)
			return
		}
	}
	response.Ok(c)
	return
}

func ScaleDeployment(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	var scaleData k8s.ScaleDeployment
	err2 := c.ShouldBindJSON(&scaleData)
	if err2 != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	err = deployment.ScaleDeployment(client, scaleData.Namespace, scaleData.DeploymentName, *scaleData.ScaleNumber)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

func DetailDeploymentController(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	namespace := parser.ParseNamespaceParameter(c)
	name := parser.ParseNameParameter(c)

	data, err := deployment.GetDeploymentDetail(client, namespace, name)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "操作成功", c)
}

func RestartDeploymentController(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	var restartDeployment k8s.RestartDeployment
	err2 := c.ShouldBindJSON(&restartDeployment)
	if err2 != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	err3 := deployment.RestartDeployment(client, restartDeployment.Namespace, restartDeployment.DeploymentName)
	if err3 != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
	}
	response.Ok(c)
	return
}

func GetDeploymentToServiceController(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}

	var Deployment k8s.RestartDeployment
	err2 := c.ShouldBindJSON(&Deployment)
	if err2 != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}

	data, err := service.GetToService(client, Deployment.Namespace, Deployment.DeploymentName)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), c)
		return
	}
	response.OkWithData(data, c)
	return
}

func RollBackDeploymentController(c *gin.Context) {
	client, err := pkg.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	var rollback k8s.RollbackDeployment
	rollbackParamsErr := c.ShouldBindJSON(&rollback)
	if err != nil {
		response.FailWithMessage(response.ParamError, rollbackParamsErr.Error(), c)
		return
	}
	rollbackErr := deployment.RollDeployment(client, rollback.DeploymentName, rollback.Namespace, *rollback.ReVersion)
	common.Log.Info(fmt.Sprintf("rollbackErr: %v", rollbackErr))
	if rollbackErr != nil {
		response.FailWithMessage(response.ERROR, err.Error(), c)
		return
	}
	response.Ok(c)

}
