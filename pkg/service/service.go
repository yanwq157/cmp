package service

import (
	"cmp/common"
	"cmp/model/k8s"
	k8scommon "cmp/pkg/common"
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type Service struct {
	ObjectMeta k8s.ObjectMeta `json:"objectMeta"`
	TypeMeta   k8s.TypeMeta   `json:"typeMeta"`

	InternalEndpoint  k8scommon.Endpoint   `json:"internalEndpoint"`
	ExternalEndpoints []k8scommon.Endpoint `json:"externalEndpoints"`
	Selector          map[string]string    `json:"selector"`
	Type              v1.ServiceType       `json:"type"`
	ClusterIP         string               `json:"clusterIP"`
}

type ServiceList struct {
	ListMeta k8s.ListMeta `json:"listMeta"`
	Services []Service    `json:"services"`
}

func ToService(service *v1.Service) Service {
	return Service{
		ObjectMeta:        k8s.NewObjectMeta(service.ObjectMeta),
		TypeMeta:          k8s.NewTypeMeta(k8s.ResourceKindDeployment),
		InternalEndpoint:  k8scommon.GetInternalEndpoint(service.Name, service.Namespace, service.Spec.Ports),
		ExternalEndpoints: k8scommon.GetExternalEndpoints(service),
		Selector:          service.Spec.Selector,
		ClusterIP:         service.Spec.ClusterIP,
		Type:              service.Spec.Type,
	}
}

func DeleteService(client *kubernetes.Clientset, ns string, serviceName string) error {
	common.Log.Info(fmt.Sprintf("请求删除Service：%v,namespace:%v", serviceName, ns))
	return client.CoreV1().Services(ns).Delete(
		context.TODO(),
		serviceName,
		metav1.DeleteOptions{},
	)
}

func GetToService(client *kubernetes.Clientset, namespace string, name string) (*ServiceList, error) {
	serviceList := &ServiceList{
		Services: make([]Service, 0),
	}
	svcList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	common.Log.Info("开始获取svc")
	if err != nil {
		return nil, err
	}
	for _, svc := range svcList.Items {
		if strings.Contains(svc.Name, name) {
			serviceList.Services = append(serviceList.Services, ToService(&svc))
			serviceList.ListMeta = k8s.ListMeta{
				TotalItems: len(serviceList.Services),
			}
			return serviceList, nil
		}
	}
	common.Log.Warn(fmt.Sprintf("没有找到所关联的svc：namespace:%s,name:%s", namespace, name))
	return nil, errors.New("没有找到所关联的svc")
}
