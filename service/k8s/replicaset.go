package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
)

type replicaSetService struct {
}

func (rss *replicaSetService) getClient(env, namespace string) v1.ReplicaSetInterface {
	return K8sClient.getClientSets(env).AppsV1().ReplicaSets(namespace)
}

func (rss replicaSetService) Get(env, namespace, replicasetName string) (*appsv1.ReplicaSet, error) {
	return rss.getClient(env, namespace).Get(replicasetName, metav1.GetOptions{})
}
