package k8s

import (
	"github.com/jhonnli/container-orchestration-api/api/k8s"
	corev1 "k8s.io/api/core/v1"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

func NewEventService() k8s.EventInterface {
	return &eventService{}
}

type eventService struct {
}

func (es *eventService) getClient(env, namespace string) v1.EventInterface {
	return K8sClient.getClientSets(env).CoreV1().Events(namespace)
}

func (es eventService) Get(env, namespace, resourceName string) ([]corev1.Event, error) {
	data, err := es.getClient(env, namespace).List(meta1.ListOptions{
		FieldSelector: "involvedObject.name=" + resourceName,
	})
	return data.Items, err
}
