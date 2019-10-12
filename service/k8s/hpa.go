package k8s

import (
	"errors"
	"fmt"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/logs"
	"k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	"strconv"
)

func NewHPAService() k8s2.HPAInterface {
	return &hpaService{}
}

type hpaService struct {
}

func (hpa *hpaService) getClient(env, namespace string) v2.HorizontalPodAutoscalerInterface {
	return K8sClient.getClientSets(env).AutoscalingV2beta1().HorizontalPodAutoscalers(namespace)
}

func (hpa *hpaService) Apply(env string, param k8s.HPAParam) (*v2beta1.HorizontalPodAutoscaler, error) {
	horizontalPodAutoscaler, err := hpa.castToHorizontalPodAutoscaler(param)
	if err != nil {
		logs.Warn("apply hap转换时发生异常,原因: %s\n", err.Error())
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	hpaCurrent, err := hpa.getClient(env, param.MetaData.Namespace).Get(param.MetaData.Name, metav1.GetOptions{})
	if err != nil || hpaCurrent.Name == "" {
		return hpa.Create(env, param)
	}
	horizontalPodAutoscaler.ObjectMeta = hpaCurrent.ObjectMeta
	data, err := hpa.getClient(env, param.MetaData.Namespace).Update(horizontalPodAutoscaler)
	if err != nil {
		logs.Error(err)
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	return data, err
}

func (hpa *hpaService) Create(env string, param k8s.HPAParam) (*v2beta1.HorizontalPodAutoscaler, error) {
	horizontalPodAutoscaler, err := hpa.castToHorizontalPodAutoscaler(param)
	if err != nil {
		logs.Warn("create hap转换时发生异常,原因: %s\n", err.Error())
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	data, err := hpa.getClient(env, param.MetaData.Namespace).Create(horizontalPodAutoscaler)

	if err != nil {
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	return data, err
}

func (hpa *hpaService) List(env string, namespace string) (*v2beta1.HorizontalPodAutoscalerList, error) {
	return hpa.getClient(env, namespace).List(metav1.ListOptions{})
}

func (hpa *hpaService) Get(env string, namespace, name string) (*v2beta1.HorizontalPodAutoscaler, error) {
	return hpa.getClient(env, namespace).Get(name, metav1.GetOptions{})
}

func (hpa *hpaService) Delete(env string, namepsace, name string) error {
	return hpa.getClient(env, namepsace).Delete(name, &metav1.DeleteOptions{})
}

func (hpa *hpaService) Update(env string, param k8s.HPAParam) (*v2beta1.HorizontalPodAutoscaler, error) {
	horizontalPodAutoscaler, err := hpa.castToHorizontalPodAutoscaler(param)
	if err != nil {
		logs.Warn("update hap转换时发生异常,原因: %s\n", err.Error())
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	hpaCurrent, err := hpa.getClient(env, param.MetaData.Namespace).Get(param.MetaData.Name, metav1.GetOptions{})
	if err != nil {
		logs.Warn("update中获取hap %s下%s的hpa时发生异常,原因: %s\n", param.MetaData.Namespace, param.MetaData.Name, err.Error())
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	horizontalPodAutoscaler.ObjectMeta = hpaCurrent.ObjectMeta
	data, err := hpa.getClient(env, param.MetaData.Namespace).Update(horizontalPodAutoscaler)
	if err != nil {
		logs.Error(err)
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	return data, err
}

func (hpa *hpaService) castMetrics(metricsParam k8s.HPAMetricsParam) ([]v2beta1.MetricSpec, error) {
	metrics := make([]v2beta1.MetricSpec, 0)
	cpuMetrics := v2beta1.MetricSpec{
		Type: v2beta1.ResourceMetricSourceType,
		Resource: &v2beta1.ResourceMetricSource{
			Name:                     "cpu",
			TargetAverageUtilization: int32ToPoint(metricsParam.Cpu),
		},
	}
	memory := strconv.Itoa(int(metricsParam.Memory)) + "Mi"
	quantity, err := resource.ParseQuantity(memory)
	if err != nil {
		msg := fmt.Sprint("解析hpa quantity[%s]失败, 原因: %s\n", quantity, err.Error())
		logs.Warn(msg)
		return metrics, errors.New(msg)
	}
	memoryMetrics := v2beta1.MetricSpec{
		Type: v2beta1.ResourceMetricSourceType,
		Resource: &v2beta1.ResourceMetricSource{
			Name:               "memory",
			TargetAverageValue: &quantity,
		},
	}
	metrics = append(metrics, cpuMetrics)
	metrics = append(metrics, memoryMetrics)
	return metrics, nil
}

func (hpa *hpaService) castToHorizontalPodAutoscaler(params k8s.HPAParam) (*v2beta1.HorizontalPodAutoscaler, error) {
	metrics, err := hpa.castMetrics(params.Spec.Metrics)
	if err != nil {
		return &v2beta1.HorizontalPodAutoscaler{}, err
	}
	spec := v2beta1.HorizontalPodAutoscalerSpec{
		MinReplicas: &params.Spec.MinReplicas,
		MaxReplicas: params.Spec.MaxReplicas,
		ScaleTargetRef: v2beta1.CrossVersionObjectReference{
			Kind:       "Deployment",
			Name:       params.Spec.ScaleTargetRefName,
			APIVersion: "extensions/v1beta1",
		},
		Metrics: metrics,
	}
	metaData := castObjectMeta(params.MetaData)
	return &v2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metaData,
		Spec:       spec,
	}, nil
}
