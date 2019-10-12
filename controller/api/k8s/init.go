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

var engine *gin.Engine

func Init(e *gin.Engine) {
	engine = e
	InitApplication()
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
