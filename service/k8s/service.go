package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/exception"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"strconv"
	"strings"
)

func NewSvcService() k8s2.ServiceInterface {
	return &svcService{}
}

type svcService struct {
}

func (svc *svcService) getClient(env, namespace string) v1.ServiceInterface {
	return K8sClient.getClientset(env).CoreV1().Services(namespace)
}

func (svc *svcService) castToService(param k8s.ServiceParam) v12.Service {
	service := v12.Service{
		ObjectMeta: svc.castObjectMeta(param.MetaData),
		Spec:       svc.castSpec(param.Spec),
	}
	return service
}

func (svc *svcService) castObjectMeta(param k8s.ObjectMetaParam) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      param.Name,
		Namespace: param.Namespace,
		Labels:    param.Labels,
	}
}

func (svc *svcService) castServicePort(param k8s.PortParam) v12.ServicePort {
	return v12.ServicePort{
		Name:     "tcp-" + strconv.Itoa(int(param.Port)),
		Protocol: v12.ProtocolTCP,
		TargetPort: intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: param.TargetPort,
		},
		Port: param.Port,
	}
}

func (svc *svcService) castSpec(param k8s.ServiceSpecParam) v12.ServiceSpec {
	ports := make([]v12.ServicePort, 0)
	for _, portParam := range param.Ports {
		port := svc.castServicePort(portParam)
		ports = append(ports, port)
	}
	return v12.ServiceSpec{
		Ports:    ports,
		Type:     svc.getServiceSpecType(param.Type),
		Selector: param.Selector,
	}
}

func (svc *svcService) getServiceSpecType(typeStr string) v12.ServiceType {
	typeStr = strings.ToLower(typeStr)
	switch typeStr {
	case "clusterip":
		return v12.ServiceTypeClusterIP
	case "nodeport":
		return v12.ServiceTypeNodePort
	case "loadbalancer":
		return v12.ServiceTypeLoadBalancer
	case "externalname":
		return v12.ServiceTypeExternalName
	default:
		return v12.ServiceTypeClusterIP
	}
}

func (svc *svcService) Get(env, namespace string, svcname string) (*v12.Service, error) {
	return svc.getClient(env, namespace).Get(svcname, metav1.GetOptions{})
}

func (svc svcService) Delete(env, namespace string, svcname string) error {
	return svc.getClient(env, namespace).Delete(svcname, &metav1.DeleteOptions{})
}

func (svc svcService) List(env, namespace string) (*v12.ServiceList, error) {
	return svc.getClient(env, namespace).List(metav1.ListOptions{})
}

func (svc *svcService) Create(env string, param k8s.ServiceParam) (*v12.Service, error) {
	client := svc.getClient(env, param.MetaData.Namespace)
	service := svc.castToService(param)
	return client.Create(&service)
}

func (svc *svcService) Update(env string, param k8s.ServiceParam) (*v12.Service, error) {
	client := svc.getClient(env, param.MetaData.Namespace)
	service := svc.castToService(param)
	exist, service1 := svc.exist(env, param.MetaData.Namespace, param.MetaData.Name)
	if exist {
		service1.Spec.Ports = service.Spec.Ports
		service1.Spec.Selector = service.Spec.Selector
		service1.Spec.Type = service.Spec.Type
		return client.Update(service1)
	} else {
		return nil, exception.NewError("service.update.not_exist", "服务不存在")
	}
}

/**
service创建或更新，如果更新时，不允许更新metadata信息
*/
func (svc *svcService) Apply(env string, param k8s.ServiceParam) (*v12.Service, error) {
	client := svc.getClient(env, param.MetaData.Namespace)
	service := svc.castToService(param)
	exist, service1 := svc.exist(env, param.MetaData.Namespace, param.MetaData.Name)
	if exist {
		service1.Spec.Ports = service.Spec.Ports
		service1.Spec.Selector = service.Spec.Selector
		service1.Spec.Type = service.Spec.Type
		return client.Update(service1)
	} else {
		return client.Create(&service)
	}
}

func (svc *svcService) exist(env, namespace, name string) (bool, *v12.Service) {
	result := true
	service, err := svc.Get(env, namespace, name)
	if err != nil {
		result = false
	} else {
		if service.Name == "" {
			result = false
		}
	}
	return result, service
}
