package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/golib/logs"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"sync"
)

var ingMutex sync.Mutex

func NewIngressService() k8s2.IngressInterface {
	return &ingressService{}
}

type ingressService struct {
}

func (ing *ingressService) getClient(env, namespace string) v1beta1.IngressInterface {
	return K8sClient.getClientset(env).ExtensionsV1beta1().Ingresses(namespace)
}

func (ing *ingressService) Apply(env string, param k8s.IngressParam) (*v1beta12.Ingress, error) {
	ingress := ing.castToIngress(param)
	ing1, err := ing.Get(env, param.MetaData.Namespace, param.MetaData.Name)
	if err != nil || ing1.Name == "" {
		logs.Info("获取%s环境%s下Ingress[%s]发生异常，%s\n", env, param.MetaData.Namespace,
			param.MetaData.Name, err.Error())
		return ing.Create(env, param)
	}
	ingress.ObjectMeta = ing1.ObjectMeta
	return ing.getClient(env, param.MetaData.Namespace).Update(ingress)
}

func (ing *ingressService) Update(env string, param k8s.IngressParam) (*v1beta12.Ingress, error) {
	ingress := ing.castToIngress(param)
	ing1, err := ing.Get(env, param.MetaData.Namespace, param.MetaData.Name)
	if err != nil {
		logs.Info("获取%s环境%s下Ingress[%s]发生异常，%s\n", env, param.MetaData.Namespace,
			param.MetaData.Name, err.Error())
		return &v1beta12.Ingress{}, err
	} else {
		ingress.ObjectMeta = ing1.ObjectMeta
	}
	return ing.getClient(env, param.MetaData.Namespace).Update(ingress)
}

func (ing *ingressService) Get(env, namespace, ingName string) (*v1beta12.Ingress, error) {
	return ing.getClient(env, namespace).Get(ingName, v1.GetOptions{})
}

func (ing *ingressService) List(env, namespace string) (*v1beta12.IngressList, error) {
	return ing.getClient(env, namespace).List(v1.ListOptions{})
}

func (ing *ingressService) Create(env string, param k8s.IngressParam) (*v1beta12.Ingress, error) {
	ingress := ing.castToIngress(param)
	return ing.getClient(env, param.MetaData.Namespace).Create(ingress)
}

func (ing *ingressService) castToIngress(param k8s.IngressParam) *v1beta12.Ingress {
	return &v1beta12.Ingress{
		ObjectMeta: castObjectMeta(param.MetaData),
		Spec:       ing.castIngressSpec(param.Spec),
	}
}

func (ing *ingressService) castIngressSpec(param k8s.IngressSpecParam) v1beta12.IngressSpec {
	return v1beta12.IngressSpec{
		TLS:   ing.castToIngressTLS(param.TLS),
		Rules: ing.castToIngressRule(param.Rules),
	}
}

func (ing *ingressService) castToIngressTLS(params []k8s.IngressTLSParam) []v1beta12.IngressTLS {
	result := make([]v1beta12.IngressTLS, 0)
	for _, param := range params {
		tls := v1beta12.IngressTLS{
			SecretName: param.SecretName,
			Hosts:      param.Hosts,
		}
		result = append(result, tls)
	}
	return result
}

func (ing *ingressService) castToIngressRule(params []k8s.RuleParam) []v1beta12.IngressRule {
	result := make([]v1beta12.IngressRule, 0)
	for _, item := range params {
		data := v1beta12.IngressRule{
			Host: item.Host,
			IngressRuleValue: v1beta12.IngressRuleValue{
				HTTP: &v1beta12.HTTPIngressRuleValue{
					Paths: ing.castToHttpIngressPath(item.Http),
				},
			},
		}
		result = append(result, data)
	}
	return result
}

func (ing *ingressService) castToHttpIngressPath(param k8s.HttpParam) []v1beta12.HTTPIngressPath {
	result := make([]v1beta12.HTTPIngressPath, 0)
	for _, item := range param.Paths {
		data := v1beta12.HTTPIngressPath{
			Path: item.Path,
			Backend: v1beta12.IngressBackend{
				ServiceName: item.Backend.ServiceName,
				ServicePort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: item.Backend.ServicePort,
				},
			},
		}
		result = append(result, data)
	}
	return result
}

func (ing *ingressService) Delete(env, namespace, ingName string) error {
	return ing.getClient(env, namespace).Delete(ingName, &v1.DeleteOptions{})
}
