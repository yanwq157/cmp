package routers

import (
	"cmp/api/v1"
	"github.com/gin-gonic/gin"
)

func InitContainerRouter(r *gin.RouterGroup) {
	K8sClusterRouter := r.Group("clusters")
	{
		K8sClusterRouter.POST("addCluster", v1.AddCluster)
	}
}
