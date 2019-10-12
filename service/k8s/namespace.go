package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

func NewNamespaceService() k8s2.NamespaceInterface {
	return &namespaceService{}
}

type namespaceService struct {
}

func (ns *namespaceService) getClient(env string) v1.NamespaceInterface {
	return K8sClient.getClientset(env).CoreV1().Namespaces()
}

func (ns *namespaceService) Create(env string, param k8s.NamespaceParam) (*v12.Namespace, error) {
	nsvo := &v12.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: param.Name,
		},
	}
	return ns.getClient(env).Create(nsvo)
}

func (ns *namespaceService) List(env string) (*v12.NamespaceList, error) {

	return ns.getClient(env).List(metav1.ListOptions{})
}

func (ns *namespaceService) Delete(env, namespace string) error {
	return ns.getClient(env).Delete(namespace, &metav1.DeleteOptions{})
}

func (ns *namespaceService) Get(env, name string) (*v12.Namespace, error) {
	return ns.getClient(env).Get(name, metav1.GetOptions{})
}
