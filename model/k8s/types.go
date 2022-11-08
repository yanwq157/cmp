package k8s

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

func NewObjectMeta(k8SObjectMeta metaV1.ObjectMeta) ObjectMeta {
	return ObjectMeta{
		Name:              k8SObjectMeta.Name,
		Namespace:         k8SObjectMeta.Namespace,
		Labels:            k8SObjectMeta.Labels,
		CreationTimestamp: Time(k8SObjectMeta.CreationTimestamp),
		Annotations:       k8SObjectMeta.Annotations,
	}
}

// NewTypeMeta creates new type mete for the resource kind.
func NewTypeMeta(kind ResourceKind) TypeMeta {
	return TypeMeta{
		Kind: kind,
	}
}

//列出所有资源而不进行任何过滤
var ListEverything = metaV1.ListOptions{
	//通过标签限制返回对象列表的选择器。默认为一切
	LabelSelector: labels.Everything().String(),
	//通过字段限制返回对象列表的选择器。默认为一切
	FieldSelector: fields.Everything().String(),
}

const (
	ResourceKindDeployment = "deployment"
	ResourceKindPod        = "pod"
)
