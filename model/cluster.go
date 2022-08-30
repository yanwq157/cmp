package model

import "gorm.io/gorm"

type Cluster struct {
	gorm.Model
	ClusterName    string `json:"clusterName" gorm:"comment:集群名称" form:"clusterName" binding:"required"`
	KubeConfig     string `json:"kubeConfig" gorm:"comment:集群凭证;type:varchar(12800)" binding:"required"`
	ClusterVersion string `json:"clusterVersion" gorm:"comment:集群版本"`
	//NodeNumber     int    `json:"nodeNumber" gorm:"comment:节点数"`
}
