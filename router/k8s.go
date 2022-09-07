package routers

import (
	"cmp/api/v1"
	"github.com/gin-gonic/gin"
)

func InitContainerRouter(r *gin.RouterGroup) {
	K8sClusterRouter := r.Group("k8s")
	{
		K8sClusterRouter.POST("cluster", v1.CreateCluster)
		K8sClusterRouter.GET("cluster", v1.ListCluster)
		K8sClusterRouter.POST("cluster/delete", v1.DelCluster)
		K8sClusterRouter.GET("cluster/detail", v1.GetK8SClusterDetail)

		K8sClusterRouter.GET("node", v1.GetNodes)

	}
}
