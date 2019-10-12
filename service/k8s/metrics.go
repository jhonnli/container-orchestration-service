package k8s

import (
	"fmt"
	"github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/metrics"
)

func NewMetricsService() k8s.MetricsInterface {
	return &metricsService{}
}

type metricsService struct {
}

func (ms metricsService) getRestClientParam() RestClientParam {
	return RestClientParam{"/apis", "metrics.k8s.io", "v1beta1"}
}

func (ms metricsService) getRequestPrefix() string {
	return "/apis/metrics.k8s.io/v1beta1"
}

func (ms metricsService) ListMetricsOfNode(env string) ([]metrics.NodeMetrics, error) {
	metricsList := &metrics.NodeMetricsList{}
	err := K8sClient.getRestClient(env, ms.getRestClientParam()).Get().
		RequestURI(ms.getRequestPrefix() + "/nodes").Do().Into(metricsList)
	return metricsList.Items, err
}

func (ms metricsService) GetMetricsOfNodeByNodeName(env, node string) (*metrics.NodeMetrics, error) {
	data := &metrics.NodeMetrics{}
	err := K8sClient.getRestClient(env, ms.getRestClientParam()).Get().
		RequestURI(ms.getRequestPrefix() + "/nodes/" + node).Do().Into(data)
	value, _ := json.Marshal(data)
	fmt.Println(string(value))
	return data, err
}

func (ms metricsService) ListMetricsOfNamespace(env, namespace string) ([]metrics.PodMetrics, error) {
	metricsList := &metrics.PodMetricsList{}
	err := K8sClient.getRestClient(env, ms.getRestClientParam()).Get().
		RequestURI(ms.getRequestPrefix() + "/namespaces/" + namespace + "/pods").Do().Into(metricsList)
	return metricsList.Items, err
}

func (ms metricsService) GetMetricsOfNamespaceByPodName(env, namespace, podName string) (*metrics.PodMetrics, error) {
	data := &metrics.PodMetrics{}
	err := K8sClient.getRestClient(env, ms.getRestClientParam()).Get().
		RequestURI(ms.getRequestPrefix() + "/namespaces/" + namespace + "/pods/" + podName).Do().Into(data)
	return data, err
}
