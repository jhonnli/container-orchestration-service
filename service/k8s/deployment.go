package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	v12 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
	"strings"
)

func NewDeploymentService() k8s2.DeploymentInterface {
	return &deploymentService{}
}

type deploymentService struct {
}

func (deploy *deploymentService) getClient(env, namespace string) v1.DeploymentInterface {
	return K8sClient.getClientset(env).AppsV1().Deployments(namespace)
}

func (ds *deploymentService) Apply(env string, param k8s.DeploymentParam) (*v12.Deployment, error) {
	deploy := ds.castToDeployment(param)
	deployCurrent, err := ds.Get(env, param.MetaData.Namespace, param.MetaData.Name)
	if err != nil || deployCurrent.Name == "" {
		return ds.Create(env, param)
	}
	deploy.ObjectMeta = deployCurrent.ObjectMeta
	return ds.getClient(env, param.MetaData.Namespace).Update(deploy)
}

func (ds *deploymentService) Create(env string, param k8s.DeploymentParam) (*v12.Deployment, error) {
	if param.MetaData.Namespace == "" {
		param.MetaData.Namespace = "default"
	}
	deployment := ds.castToDeployment(param)

	return ds.getClient(env, param.MetaData.Namespace).Create(deployment)
}

func (ds *deploymentService) castToDeploymentSpec(meta metav1.ObjectMeta, param k8s.DeploymentSpecParam) v12.DeploymentSpec {
	selector := make(map[string]string)
	selector["app"] = meta.Name
	var labels map[string]string
	if meta.Labels == nil {
		labels = make(map[string]string)
	}
	labels["app"] = meta.Name
	meta.Labels = labels
	return v12.DeploymentSpec{
		Replicas:                &param.Replicas,
		ProgressDeadlineSeconds: int32ToPoint(600),
		RevisionHistoryLimit:    int32ToPoint(2),
		Selector: &metav1.LabelSelector{
			MatchLabels: selector,
		},
		Strategy: v12.DeploymentStrategy{
			RollingUpdate: &v12.RollingUpdateDeployment{
				MaxSurge: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "25%",
				},
				MaxUnavailable: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "25%",
				},
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: meta,
			Spec: corev1.PodSpec{
				Containers:                    ds.castToContainers(meta.Name, param.Template.Spec.Containers),
				RestartPolicy:                 corev1.RestartPolicyAlways,
				TerminationGracePeriodSeconds: int64ToPoint(30),
				Volumes:                       ds.castToVolumes(),
				NodeSelector:                  param.Template.Spec.NodeSelector,
			},
		},
	}
}

func (ds *deploymentService) castToVolumes() []corev1.Volume {
	result := make([]corev1.Volume, 0)
	result = append(result, corev1.Volume{
		Name: "localtime",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/etc/localtime",
			},
		},
	})
	return result
}

func (ds *deploymentService) castToContainers(name string, params []k8s.ContainerParam) []corev1.Container {
	result := make([]corev1.Container, 0)
	for _, item := range params {
		ct := corev1.Container{
			Name:  name,
			Image: item.Image,
			LivenessProbe: &corev1.Probe{
				FailureThreshold:    5,
				InitialDelaySeconds: 30,
				PeriodSeconds:       60,
				SuccessThreshold:    1,
				TimeoutSeconds:      2,
				Handler: corev1.Handler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: item.ServicePort,
						},
					},
				},
			},
			ReadinessProbe: &corev1.Probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       30,
				SuccessThreshold:    1,
				TimeoutSeconds:      10,
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: item.ServicePort,
						},
						Path: item.HealthCheckPath,
					},
				},
			},
			Resources:    ds.castToContainerResource(item.Resources),
			VolumeMounts: ds.castToVolumeMounts(item.VolumeMounts),
		}
		result = append(result, ct)
	}
	return result
}

func (ds *deploymentService) castToContainerResource(param k8s.ContainerResourceParams) corev1.ResourceRequirements {
	if param.Limits.Cpu == "" {
		param.Limits.Cpu = param.Requests.Cpu
	}
	if param.Limits.Memory == "" {
		param.Limits.Memory = param.Requests.Memory
	}

	return corev1.ResourceRequirements{
		Limits:   ds.castToContainerResourceQuota(param.Limits),
		Requests: ds.castToContainerResourceQuota(param.Requests),
	}
}

func (ds *deploymentService) castToContainerResourceQuota(param k8s.ContainerResourceParam) corev1.ResourceList {
	result := make(map[corev1.ResourceName]resource.Quantity)
	cpuQuantity, _ := resource.ParseQuantity(param.Cpu)
	result[corev1.ResourceCPU] = cpuQuantity

	if !strings.HasSuffix(param.Memory, "Mi") {
		param.Memory = param.Memory + "Mi"
	}
	memoryQuantity, _ := resource.ParseQuantity(param.Memory)
	result[corev1.ResourceMemory] = memoryQuantity
	return result
}

func (ds *deploymentService) castToVolumeMounts(params []k8s.VolumeMountParam) []corev1.VolumeMount {
	result := make([]corev1.VolumeMount, 0)
	for _, item := range params {
		result = append(result, ds.castToVolume(item))
	}

	result = append(result, ds.castToVolume(k8s.VolumeMountParam{
		MountPath: "/etc/localtime",
		ReadOnly:  true,
		Name:      "localtime",
	}))
	return result
}

func (ds *deploymentService) castToVolume(param k8s.VolumeMountParam) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      param.Name,
		ReadOnly:  param.ReadOnly,
		MountPath: param.MountPath,
	}
}

func (ds *deploymentService) castToDeployment(param k8s.DeploymentParam) *v12.Deployment {
	metaData := castObjectMeta(param.MetaData)
	deployment := &v12.Deployment{
		ObjectMeta: metaData,
		Spec:       ds.castToDeploymentSpec(metaData, param.Spec),
	}
	return deployment
}

func (ds *deploymentService) Update(env string, param k8s.DeploymentParam) (*v12.Deployment, error) {
	deploy := ds.castToDeployment(param)
	deployCurrent, err := ds.Get(env, param.MetaData.Namespace, param.MetaData.Name)
	if err != nil {
		return &v12.Deployment{}, err
	}
	deploy.ObjectMeta = deployCurrent.ObjectMeta
	return ds.getClient(env, param.MetaData.Namespace).Update(deploy)
}

func (ds *deploymentService) Get(env, namespaceName, deployName string) (*v12.Deployment, error) {
	return ds.getClient(env, namespaceName).Get(deployName, metav1.GetOptions{})
}

func (ds *deploymentService) List(env, namespaceName string) (*v12.DeploymentList, error) {
	return ds.getClient(env, namespaceName).List(metav1.ListOptions{})
}

func (ds *deploymentService) Delete(env, namespaceName, deployName string) error {
	return ds.getClient(env, namespaceName).Delete(deployName, &metav1.DeleteOptions{})
}
