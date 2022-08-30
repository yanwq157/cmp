package service

import (
	"cmp/common"
	"cmp/model"
)

func CreateCluster(cluster model.Cluster) (err error) {
	err = common.Db.Create(&cluster).Error
	return
}
