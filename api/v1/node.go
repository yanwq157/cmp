package v1

import (
	"cmp/api/v1/response"
	"cmp/pkg/init"
	"github.com/gin-gonic/gin"
)

func GetNodes(c *gin.Context) {
	client, err := init.GetClusterId(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

}
