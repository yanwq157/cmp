package v1

import (
	"cmp/model"
	"cmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddCluster(c *gin.Context) {
	d := model.Cluster{}
	var msg string
	err := c.ShouldBindJSON(&d)
	if err != nil {
		return
	}
	if err := service.CreateCluster(d); err != nil {
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			msg: d,
		})
		return
	}
}
