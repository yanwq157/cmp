package k8s

import "time"

type ObjectMeta struct {
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	CreationTimestamp Time              `json:"creationTimestamp,omitempty"`
}

type TypeMeta struct {
	Kind ResourceKind `json:"kind,omitempty"`
}

type NodeAllocatedResources struct {
	CPURequests            int64   `json:"cpuRequests"`
	CPURequestsFraction    float64 `json:"cpuRequestsFraction"`
	CPULimits              int64   `json:"cpuLimits"`
	CPULimitsFraction      float64 `json:"cpuLimitsFraction"`
	CPUCapacity            int64   `json:"cpuCapacity"`
	MemoryRequests         int64   `json:"memoryRequests"`
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`
	MemoryLimits           int64   `json:"memoryLimits"`
	MemoryLimitsFraction   float64 `json:"memoryLimitsFraction"`
	MemoryCapacity         int64   `json:"memoryCapacity"`
	AllocatedPods          int     `json:"allocatedPods"`
	PodCapacity            int64   `json:"podCapacity"`
	PodFraction            float64 `json:"podFraction"`
}

type Time struct {
	time.Time `protobuf:"-"`
}

type ResourceKind string

type Unschedulable bool

type NodeIP string

type UID string
