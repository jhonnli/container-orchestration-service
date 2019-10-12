package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/service/util"
	coreV1 "k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/policy/v1beta1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
)

func NewPodService() k8s2.PodInterface {
	return &podService{}
}

type podService struct {
}

func (ps *podService) getClient(env, namespace string) v1.PodInterface {
	return K8sClient.getClientset(env).CoreV1().Pods(namespace)
}

func (ps *podService) ListByNamespace(envName, namespace string) (*coreV1.PodList, error) {

	return ps.getClient(envName, namespace).List(v12.ListOptions{})
}

func (ds podService) ListPodByDeployment(env, nsname, deployName string) (*coreV1.PodList, error) {
	options := k8s.ListOptions{
		LabelSelector: "app=" + deployName,
	}
	return ds.listByNamespaceAndOptions(env, nsname, options)
}

func (ps podService) DeletePodByDeployment(env, nsname, deployName string) error {
	options := k8s.ListOptions{
		LabelSelector: "app=" + deployName,
	}
	listOptions := util.CastToK8sListOptions(options)
	return ps.getClient(env, nsname).DeleteCollection(&v12.DeleteOptions{}, listOptions)
}

func (ps *podService) listByNamespaceAndOptions(envName, namespace string, options k8s.ListOptions) (*coreV1.PodList, error) {
	listOptions := util.CastToK8sListOptions(options)
	return ps.getClient(envName, namespace).List(listOptions)
}

func (ps *podService) ListByNode(env, nodename string) (*coreV1.PodList, error) {
	param := RestClientParam{"/api", "core", "v1"}
	client := K8sClient.getRestClient(env, param)
	data := &coreV1.PodList{}
	err := client.Get().RequestURI("/api/v1/pods").
		Param("fieldSelector", "spec.nodeName="+nodename).
		Do().Into(data)
	return data, err
}

func (ps *podService) GetByName(env, namespace, podName string) (*coreV1.Pod, error) {
	return ps.getClient(env, namespace).Get(podName, v12.GetOptions{})
}

func (ps *podService) Eviction(env, namespace, podName string) error {
	eviction := &v1beta12.Eviction{
		ObjectMeta: v12.ObjectMeta{
			Namespace: namespace,
			Name:      podName,
		},
		DeleteOptions: &v12.DeleteOptions{},
	}
	return ps.getClient(env, namespace).Evict(eviction)
}

func (ps podService) EvictionPods(env, namespace string, data []string) (bool, error) {
	result := true
	var err error
	for _, item := range data {
		eviction := &v1beta12.Eviction{
			ObjectMeta: v12.ObjectMeta{
				Namespace: namespace,
				Name:      item,
			},
			DeleteOptions: &v12.DeleteOptions{},
		}
		err = ps.getClient(env, namespace).Evict(eviction)
		if err != nil {
			result = false
			break
		}
	}
	return result, err
}

func (ps *podService) HealthCheck(env, namespace, podPrefix string) (k8s.PodHealthResponse, error) {
	data, err := ps.getClient(env, namespace).List(v12.ListOptions{})
	podList := make([]k8s.PodHealth, 0)
	if err != nil {
		return k8s.PodHealthResponse{Check: false}, err
	}
	var exist bool = false
	for _, item := range data.Items {
		if strings.HasPrefix(item.Name, podPrefix) {
			exist = true
			if string(item.Status.Phase) != "Running" {
				podList = append(podList, CastPodToPodHealth(item))
			} else {
				for _, containerStatus := range item.Status.ContainerStatuses {
					if !containerStatus.Ready {
						podList = append(podList, CastPodToPodHealth(item))
					}
				}
			}
		}
	}
	if !exist {
		return k8s.PodHealthResponse{Check: false}, nil
	}
	check := false
	if len(podList) == 0 {
		check = true
	}
	return k8s.PodHealthResponse{
		Check:  check,
		Detail: podList,
	}, err
}

func (ps podService) Log(env, namespace, podName string, logOptionParam k8s.PodLogOptionParam) (string, error) {
	logOption := ps.castToPodLogOptions(logOptionParam)

	client := ps.getClient(env, namespace).GetLogs(podName, logOption)
	result, err := client.Do().Raw()
	return string(result), err
}

func (ps *podService) castToPodLogOptions(logOptionParam k8s.PodLogOptionParam) *coreV1.PodLogOptions {
	logOption := &coreV1.PodLogOptions{}
	if logOptionParam.Container != "" {
		logOption.Container = logOptionParam.Container
	}
	if logOptionParam.LimitBytes != 0 {
		logOption.LimitBytes = &logOptionParam.LimitBytes
	}
	if logOptionParam.Previous {
		logOption.Previous = logOptionParam.Previous
	}
	if logOptionParam.SinceSeconds != 0 {
		logOption.SinceSeconds = &logOptionParam.SinceSeconds
	}
	if logOptionParam.TailLines != 0 {
		logOption.TailLines = &logOptionParam.TailLines
	}
	if logOptionParam.Timestamps {
		logOption.Timestamps = logOptionParam.Timestamps
	}
	return logOption
}
