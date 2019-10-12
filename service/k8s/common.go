package k8s

import (
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/json-iterator/go"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func castObjectMeta(param k8s.ObjectMetaParam) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      param.Name,
		Namespace: param.Namespace,
		Labels:    param.Labels,
	}
}

func castK8sObjectMeta(param metav1.ObjectMeta) k8s.ObjectMeta {
	return k8s.ObjectMeta{
		Name:      param.Name,
		Namespace: param.Namespace,
		Labels:    param.Labels,
	}
}

func int32ToPoint(data int32) *int32 {
	return &data
}

func int64ToPoint(data int64) *int64 {
	return &data
}

func CastPodToPodHealth(pod v1.Pod) k8s.PodHealth {
	return k8s.PodHealth{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Phase:     string(pod.Status.Phase),
		Reason:    pod.Status.Reason,
	}
}
