package common

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Event struct {
	ObjectMeta         v1.ObjectMeta `json:"objectMeta"`
	TypeMeta           v1.TypeMeta   `json:"typeMeta"`
	Message            string        `json:"message"`
	SourceComponent    string        `json:"sourceComponent"`
	SourceHost         string        `json:"sourceHost"`
	SubObject          string        `json:"subObject"`
	SubObjectKind      string        `json:"subObjectKind"`
	SubObjectNmae      string        `json:"subObjectNmae"`
	SubobjectNamespace string        `json:"subobjectNamespace"`
	Count              int32         `json:"count"`
	FirstSeen          v1.Time       `json:"firstSeen"`
	LastSeen           v1.Time       `json:"lastSeen"`
	Reason             string        `json:"reason"`
	Type               string        `json:"type"`
}
