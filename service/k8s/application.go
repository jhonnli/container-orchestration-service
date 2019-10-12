package k8s

import (
	sapi "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/golib/logs"
)

func NewApplicationService() sapi.ApplicationInterface {
	return &applicationService{
		deploymentService: NewDeploymentService(),
		hpaService:        NewHPAService(),
		svcService:        NewSvcService(),
		ingressService:    NewIngressService(),
	}
}

type applicationService struct {
	deploymentService sapi.DeploymentInterface
	hpaService        sapi.HPAInterface
	svcService        sapi.ServiceInterface
	ingressService    sapi.IngressInterface
}

func (as *applicationService) Apply(env string, param k8s.ApplicationParam) map[string]bool {
	result := make(map[string]bool)
	if param.Deployment.MetaData.Name != "" {
		_, err := as.deploymentService.Apply(env, param.Deployment)
		if err != nil {
			logs.Info("apply deployment: 环境[%s], 参数: %s 发生错误，原因: %s\n", env, param.Deployment, err.Error())
			result["deployment"] = false
		} else {
			result["deployment"] = true
		}
	}

	if param.HorizontalPodAutoscaler.MetaData.Name != "" {
		_, err := as.hpaService.Apply(env, param.HorizontalPodAutoscaler)
		if err != nil {
			logs.Info("apply HPA: 环境[%s], 参数: %s 发生错误，原因: %s\n", env, param.HorizontalPodAutoscaler, err.Error())
			result["horizontalPodAutoscaler"] = false
		} else {
			result["horizontalPodAutoscaler"] = true
		}
	}
	if param.Service.MetaData.Name != "" {
		_, err := as.svcService.Apply(env, param.Service)
		if err != nil {
			logs.Info("apply Service: 环境[%s], 参数: %s 发生错误，原因: %s\n", env, param.Service, err.Error())
			result["service"] = false
		} else {
			result["service"] = true
		}
	}
	if param.Ingress.MetaData.Name != "" {
		_, err := as.ingressService.Apply(env, param.Ingress)
		if err != nil {
			logs.Info("apply Ingress: 环境[%s], 参数: %s 发生错误，原因: %s\n", env, param.Ingress, err.Error())
			result["ingress"] = false
		} else {
			result["ingress"] = true
		}
	}

	return result
}
