package k8s

import "gorm.io/gorm"

type Cluster struct {
	gorm.Model
	ClusterName          string `json:"clusterName" gorm:"comment:集群名称" form:"clusterName" binding:"required"`
	ClusterStatus        bool   `json:"clusterStatus" gorm:"comment:集群状态"`
	ClusterVersion       string `json:"clusterVersion" gorm:"comment:集群版本"`
	NodeNumber           int    `json:"nodeNumber" gorm:"comment:节点数"`
	ConfigFileContentStr string `json:"configFileContentStr" gorm:"comment:集群凭证;type:varchar(12800)" binding:"required"`
}

type ClusterIds struct {
	Data interface{} `json:"clusterIds"`
}
type ClusterStatus struct {
	TotalNodeNum      int     `json:"totalNodeNum"`
	ReadyNodeNum      int     `json:"readyNodeNum"`
	UnReadyNodeNum    int     `json:"unreadyNodeNum"`
	CPUAllocatable    float64 `json:"cpuAllocatable"`
	CPURequested      float64 `json:"cpuRequested"`
	MemoryAllocatable float64 `json:"memoryAllocatable"`
	MemoryRequested   float64 `json:"memoryRequested"`
}

type PaginationQ struct {
	Size    int    `form:"size" json:"size"`
	Page    int    `form:"page" json:"page"`
	Total   int64  `json:"total"`
	Keyword string `form:"keyword" json:"keyword"`
}
