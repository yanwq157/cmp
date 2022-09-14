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
		K8sClusterRouter.GET("node/detail", v1.GetNodeDetail)
		K8sClusterRouter.POST("node/schedule", v1.NodeUnschedulable)
		K8sClusterRouter.GET("node/cordon", v1.CordonNode)
		K8sClusterRouter.POST("node/collectionSchedule", v1.CollectionNodeUnschedule)
		K8sClusterRouter.POST("node/collectionCordon", v1.CollectionCordonNode)
		K8sClusterRouter.DELETE("node", v1.RemoveNode)

	}
}
