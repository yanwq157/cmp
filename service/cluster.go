package service

import (
	"cmp/common"
	"cmp/model/k8s"
	"gorm.io/gorm"
)

func CreateCluster(cluster k8s.Cluster) (err error) {
	err = common.Db.Create(&cluster).Error
	return
}
func ListCluster(p *k8s.PaginationQ, k *[]k8s.Cluster) (err error) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	offset := p.Size * (p.Page - 1)
	var total int64
	err = common.Db.Where("cluster_name like ?", "%"+p.Keyword+"%").Limit(p.Size).Offset(offset).Find(&k).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return nil
}

func DelCluster(ids k8s.ClusterIds) (err error) {
	var k k8s.Cluster
	err2 := common.Db.Delete(&k, ids.Data)
	if err2.Error != nil {
		return err2.Error
	}
	return nil
}
func GetK8sClusterID(id uint) (K8sCluster k8s.Cluster, err error) {
	err = common.Db.Where("id = ?", id).First(&K8sCluster).Error
	if err != nil {
		return K8sCluster, err
	}
	return K8sCluster, nil
}
