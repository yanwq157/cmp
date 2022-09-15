package common

import (
	"cmp/model/k8s"
	"context"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
)

type ResourceChannels struct {
	ReplicaSetList ReplicaSetListChannel
	DeploymentList DeploymentListChannel
	PodList        PodListChannel
	EventList      EventListChannel
}
type EventListChannel struct {
	List  chan *v1.EventList
	Error chan error
}

func GetEventListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) EventListChannel {
	return GetEventListChannelWithOptions(client, nsQuery, k8s.ListEverything, numReads)
}
func GetEventListChannelWithOptions(client client.Interface, nsQuery *NamespaceQuery, options metav1.ListOptions, numReads int) EventListChannel {
	channel := EventListChannel{
		List:  make(chan *v1.EventList, numReads),
		Error: make(chan error, numReads),
	}
	go func() {
		list, err := client.CoreV1().Events(nsQuery.ToRequestParam()).List(context.TODO(), options)
		var filteredItems []v1.Event
		for _, item := range list.Items {
			if nsQuery.Matcher(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}

type PodListChannel struct {
	List  chan *v1.PodList
	Error chan error
}

func GetPodListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) PodListChannel {
	return GetPodListChannelWithOptions(client, nsQuery, k8s.ListEverything, numReads)
}
func GetPodListChannelWithOptions(client client.Interface, nsQuery *NamespaceQuery, options metav1.ListOptions, numReads int) PodListChannel {
	channel := PodListChannel{
		List:  make(chan *v1.PodList, numReads),
		Error: make(chan error, numReads),
	}
	go func() {
		list, err := client.CoreV1().Pods(nsQuery.ToRequestParam()).List(context.TODO(), options)
		var filteredItems []v1.Pod
		for _, item := range list.Items {
			if nsQuery.Matcher(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}

type DeploymentListChannel struct {
	List  chan *apps.DeploymentList
	Error chan error
}

func GetDeploymentListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) DeploymentListChannel {

	channel := DeploymentListChannel{
		List:  make(chan *apps.DeploymentList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.AppsV1().Deployments(nsQuery.ToRequestParam()).List(context.TODO(), k8s.ListEverything)
		var filteredItems []apps.Deployment
		for _, item := range list.Items {
			filteredItems = append(filteredItems, item)
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}

type ReplicaSetListChannel struct {
	List  chan *apps.ReplicaSetList
	Error chan error
}

func GetReplicaSetListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) ReplicaSetListChannel {
	return GetReplicaSetListChannelWithOptions(client, nsQuery, k8s.ListEverything, numReads)
}
func GetReplicaSetListChannelWithOptions(client client.Interface, nsQuery *NamespaceQuery, options metav1.ListOptions, numReads int) ReplicaSetListChannel {
	channel := ReplicaSetListChannel{
		List:  make(chan *apps.ReplicaSetList, numReads),
		Error: make(chan error, numReads),
	}
	go func() {
		list, err := client.AppsV1().ReplicaSets(nsQuery.ToRequestParam()).List(context.TODO(), options)
		var filteredItems []apps.ReplicaSet
		for _, item := range list.Items {
			if nsQuery.Matcher(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}
