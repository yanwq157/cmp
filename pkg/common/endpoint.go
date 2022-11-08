package common

import (
	"bytes"
	api "k8s.io/api/core/v1"
)

type Endpoint struct {
	Host  string        `json:"host"`
	Ports []ServicePort `json:"ports"`
}

func GetInternalEndpoint(serviceName, namespace string, ports []api.ServicePort) Endpoint {
	name := serviceName
	if namespace != api.NamespaceDefault && len(namespace) > 0 && len(serviceName) > 0 {
		bufferName := bytes.NewBufferString(name)
		bufferName.WriteString(".")
		bufferName.WriteString(namespace)
		name = bufferName.String()
	}
	return Endpoint{
		Host:  name,
		Ports: GetServicePorts(ports),
	}
}
func GetExternalEndpoints(service *api.Service) []Endpoint {
	externalEndpoints := make([]Endpoint, 0)
	if service.Spec.Type == api.ServiceTypeLoadBalancer {
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			externalEndpoints = append(externalEndpoints, getExternalEndpoint(ingress, service.Spec.Ports))
		}
	}

	for _, ip := range service.Spec.ExternalIPs {
		externalEndpoints = append(externalEndpoints, Endpoint{
			Host:  ip,
			Ports: GetServicePorts(service.Spec.Ports),
		})
	}
	return externalEndpoints
}

func getExternalEndpoint(ingress api.LoadBalancerIngress, ports []api.ServicePort) Endpoint {
	var host string
	if ingress.Hostname != "" {
		host = ingress.Hostname
	} else {
		host = ingress.IP
	}
	return Endpoint{
		Host:  host,
		Ports: GetServicePorts(ports),
	}
}
