package common

import api "k8s.io/api/core/v1"

type NamespaceQuery struct {
	Namespace []string
}

// 返回k8sapi命名空间的查询，获取此命名空间的对象列表。如果选着了一个，则查询单个，如果未选择，则查询所有
func (n *NamespaceQuery) ToRequestParam() string {
	if len(n.Namespace) == 1 {
		return n.Namespace[0]
	}
	return api.NamespaceAll
}

//传的namespaces和查询匹配返回true
func (n *NamespaceQuery) Matches(namespaces string) bool {
	if len(n.Namespace) == 0 {
		return true
	}
	for _, queryNamespace := range n.Namespace {
		if namespaces == queryNamespace {
			return true
		}
	}
	return false
}

//查询单个命名空间的新命名空间查询
func NewSameNamespaceQuery(namespace string) *NamespaceQuery {
	return &NamespaceQuery{[]string{namespace}}
}

//查询指定的命名空间，一个或多个
func NewNamespaceQuery(namespace []string) *NamespaceQuery {
	return &NamespaceQuery{namespace}
}
