package k8s

import (
	"github.com/gin-gonic/gin"
)

const NSNAME string = "nsname"
const INGNAME string = "ingName"
const SERVICENAME string = "svcName"
const DEPLOYMENT string = "deployName"
const HPANAME string = "hpaName"
const NODENAME string = "nodeName"
const PODNAME string = "podName"
const RESOURCENAME string = "resourceName"
const CONFIGMAPNAME = "configMapName"
const SECRETNAME = "secretName"

var engine *gin.Engine

func Init(e *gin.Engine) []*gin.RouterGroup {
	data := make([]*gin.RouterGroup, 0)
	engine = e
	data = append(data, InitNamespace())
	data = append(data, InitService())
	data = append(data, InitHPA())
	data = append(data, InitDeploy())
	data = append(data, InitIngress())
	data = append(data, InitNodes())
	data = append(data, InitPodOfNode())
	data = append(data, InitPodOfNamespace())
	data = append(data, InitPodOfDeployment())
	data = append(data, InitEvent())
	data = append(data, InitMetricsOfNode())
	data = append(data, InitMetricsOfNamespace())
	data = append(data, InitConfigMap())
	data = append(data, InitSecret())

	return data
}

func getNamespaceNameFromPath(ctx *gin.Context) string {
	return ctx.Param(NSNAME)
}
func getPodNameFromPath(ctx *gin.Context) string {
	return ctx.Param(PODNAME)
}

func getIngressNameFromPath(ctx *gin.Context) string {
	return ctx.Param(INGNAME)
}

func getSerciceNameFromPath(ctx *gin.Context) string {
	return ctx.Param(SERVICENAME)
}

func getNodeNameFromPath(ctx *gin.Context) string {
	return ctx.Param(NODENAME)
}

func getHpaNameFromPath(ctx *gin.Context) string {
	return ctx.Param(HPANAME)
}

func getDeploymentNameFromPath(ctx *gin.Context) string {
	return ctx.Param(DEPLOYMENT)
}

func getConfigMapNameFromPath(ctx *gin.Context) string {
	return ctx.Param(CONFIGMAPNAME)
}

func getSecretNameFromPath(ctx *gin.Context) string {
	return ctx.Param(SECRETNAME)
}
