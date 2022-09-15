package parser

import (
	"cmp/pkg/common"
	"github.com/gin-gonic/gin"
	"strings"
)

//解析路径参数中的命名空间
//命名空间是一个逗号分隔的命名空间列表。
//没有命名空间意味着“查看所有用户命名空间”，即除了 kube-system 之外的所有内容
func ParseNamespacePathParameter(request *gin.Context) *common.NamespaceQuery {
	namespace := request.Query("namespace")
	namespaces := strings.Split(namespace, ",")
	var nonEmptyNamespaces []string
	//遍历非空的命名空间到该切片
	for _, n := range namespaces {
		if len(n) > 0 {
			nonEmptyNamespaces = append(nonEmptyNamespaces, n)
		}
	}
	return common.NewNamespaceQuery(nonEmptyNamespaces)
}

//从URL解析命名空间
func ParseNamespaceParameter(request *gin.Context) string {
	return request.Query("namespace")
}

//从URL解析name参数
func ParseNameParameter(request *gin.Context) string {
	return request.Query("name")
}
