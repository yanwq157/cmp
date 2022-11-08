package common

import api "k8s.io/api/core/v1"

type ServicePort struct {
	Port     int32        `json:"port"`
	Portocol api.Protocol `json:"portocol"`
	NodePort int32        `json:"nodePort"`
}

func GetServicePorts(apiPorts []api.ServicePort) []ServicePort {
	var ports []ServicePort
	for _, port := range apiPorts {
		ports = append(ports, ServicePort{port.Port, port.Protocol, port.NodePort})
	}
	return ports
}
