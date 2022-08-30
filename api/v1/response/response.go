package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errMsg"`
}

const (
	SUCCESS               = 0
	ERROR                 = 1
	CreateK8SClusterError = 7001
	InternalServerError   = http.StatusInternalServerError
)

const (
	OkMsg                    = "操作成功"
	NotOkMsg                 = "操作失败"
	InternalServerErrorMsg   = "服务器内部错误"
	CreateK8SClusterErrorMsg = "创建K8S集群失败"
)

var CustomError = map[int]string{
	SUCCESS:               OkMsg,
	ERROR:                 NotOkMsg,
	InternalServerError:   InternalServerErrorMsg,
	CreateK8SClusterError: CreateK8SClusterErrorMsg,
}

func ResultFail(code int, data interface{}, msg string, c *gin.Context) {
	if msg == "" {
		c.JSON(http.StatusOK, Response{
			Code:   code,
			Data:   data,
			ErrMsg: CustomError[code],
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code:   code,
			Data:   data,
			ErrMsg: msg,
		})
	}
}

func ResultOk(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}
func Ok(c *gin.Context) {
	ResultOk(SUCCESS, map[string]interface{}{}, "操作成功", c)
}
func OkWithMessage(message string, c *gin.Context) {
	ResultOk(SUCCESS, map[string]interface{}{}, message, c)
}
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	ResultOk(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	ResultFail(ERROR, map[string]interface{}{}, "操作失败", c)
}
func FailWithMessage(code int, message string, c *gin.Context) {
	ResultFail(code, map[string]interface{}{}, message, c)
}
