package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

const podPathPrefixOfNode = "/v1/envs/:env/nodes/:nodeName/pods"

var podService k8s2.PodInterface

func InitPodOfNode() *gin.RouterGroup {
	podService = k8s3.NewPodService()

	podOfNodeApi := engine.Group(podPathPrefixOfNode)
	common.AddFilter(podOfNodeApi)
	podOfNodeApi.GET("", listPodByNode)

	return podOfNodeApi
}

func InitPodOfNamespace() *gin.RouterGroup {
	podOfNSApi := engine.Group("/v1/envs/:env/namespaces/:nsname/pods")
	common.AddFilter(podOfNSApi)
	podOfNSApi.GET("", listPodByNamespace)
	podOfNSApi.GET("/:podName", getPodByName)
	podOfNSApi.GET("/:podName/logs", getPodLog)
	podOfNSApi.GET("/:podName/health", podHealthCheck)
	podOfNSApi.POST("/:podName/eviction", evictPod)

	return podOfNSApi
}

func InitPodOfDeployment() *gin.RouterGroup {
	deployGroup := engine.Group("/v1/envs/:env/namespaces/:nsname/deployments/:deployName")
	common.AddFilter(deployGroup)
	deployGroup.GET("/pods", listPodByDeployment)
	deployGroup.DELETE("/pods", deletePodByDeployment)

	return deployGroup
}

func deletePodByDeployment(context *gin.Context) {
	nsname := getNamespaceNameFromPath(context)
	deployName := getDeploymentNameFromPath(context)
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, "pod.deletePodByDeployment.param_error", env, deployName, nsname) {
		return
	}
	err := podService.DeletePodByDeployment(env, nsname, deployName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.deletePodByDeployment.error", err.Error()))
		return
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func listPodByDeployment(context *gin.Context) {
	nsname := getNamespaceNameFromPath(context)
	deployName := getDeploymentNameFromPath(context)
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, "pod.listPodByDeployment.param_error", env, deployName, nsname) {
		return
	}
	data, err := podService.ListPodByDeployment(env, nsname, deployName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.listPodByDeployment.error", err.Error()))
		return
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func getPodByName(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	podName := getPodNameFromPath(context)
	if common.ParamIsEmpty(context, env, nsname, podName) {
		return
	}
	data, err := podService.GetByName(env, nsname, podName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.get.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func getPodLog(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	podName := getPodNameFromPath(context)
	if common.ParamIsEmpty(context, env, nsname, podName) {
		return
	}
	optionParam := k8s.PodLogOptionParam{}
	err := context.BindQuery(&optionParam)
	if err != nil {
		context.JSON(http.StatusOK, common.GenParamErrorResult(err.Error()))
		return
	}
	status, err := podService.Log(env, nsname, podName, optionParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.log.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(status))
	}
}

func podHealthCheck(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	podPrefix := context.Param(PODNAME)
	if common.ParamIsEmpty(context, env, nsname, podPrefix) {
		return
	}
	status, err := podService.HealthCheck(env, nsname, podPrefix)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.healthcheck.error", err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.GenSuccessResult(status))
}

func evictPod(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	podName := getPodNameFromPath(context)
	if common.ParamIsEmpty(context, env, nsname, podName) {
		return
	}
	err := podService.Eviction(env, nsname, podName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.evictPod.error", err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.GenBoolResult())
}

func listPodByNode(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	node := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "pod.listPodByNode.param_error", env, node) {
		return
	}
	data, err := podService.ListByNode(env, node)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("pod.listbynamespace.error", err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.GenSuccessResult(data))
}

func listPodByNamespace(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	if common.ParamIsEmpty(context, "pod.listPodByNamespace.param_error", env, nsname) {
		return
	}
	data, err := podService.ListByNamespace(env, nsname)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("pod.listbynamespace.error", err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.GenSuccessResult(data))
}
