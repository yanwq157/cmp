package node

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

type Node struct {
}

func GetNodeList(client *kubernetes.Clientset, q *gin.Context) (*NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List()
}
