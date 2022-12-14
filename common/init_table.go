package common

import (
	"cmp/model/k8s"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

func MysqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		k8s.Cluster{},
	)
	if err != nil {
		Log.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	Log.Info("register table success")
}
