package k8s

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s2 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

var metricsService k8s.MetricsInterface

func InitMetricsOfNode() *gin.RouterGroup {
	metricsService = k8s2.NewMetricsService()

	metricsNodeApi := engine.Group("/v1/envs/:env/metrics/nodes")
	common.AddFilter(metricsNodeApi)
	metricsNodeApi.GET("", listMetricsOfNode)
	metricsNodeApi.GET("/:nodeName", getMetricsOfNodeByNodeName)
	return metricsNodeApi
}

func InitMetricsOfNamespace() *gin.RouterGroup {
	metricsPodApi := engine.Group("/v1/envs/:env/metrics/namespaces/:nsname/pods")
	common.AddFilter(metricsPodApi)
	metricsPodApi.GET("", listMetricsOfNamespace)
	metricsPodApi.GET("/:podName", getMetricsOfNamespaceByPodName)
	return metricsPodApi
}

func listMetricsOfNode(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, "metrics.listMetricsOfNode.param_error", env) {
		return
	}
	data, err := metricsService.ListMetricsOfNode(env)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			common.GenFailureResult("metrics.listMetricsOfNode.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func getMetricsOfNodeByNodeName(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "metrics.getMetricsOfNodeByNodeName.param_error", env, nodeName) {
		return
	}
	data, err := metricsService.GetMetricsOfNodeByNodeName(env, nodeName)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			common.GenFailureResult("metrics.getMetricsOfNodeByNodeName.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func listMetricsOfNamespace(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	namespace := getNamespaceNameFromPath(context)
	if common.ParamIsEmpty(context, "metrics.listMetricsOfNamespace.param_error", env, namespace) {
		return
	}
	data, err := metricsService.ListMetricsOfNamespace(env, namespace)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			common.GenFailureResult("metrics.listMetricsOfNamespace.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func getMetricsOfNamespaceByPodName(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	namespace := getNamespaceNameFromPath(context)
	podName := getPodNameFromPath(context)
	if common.ParamIsEmpty(context, "metrics.getMetricsOfNamespaceByPodName.param_error", env, namespace, podName) {
		return
	}
	data, err := metricsService.GetMetricsOfNamespaceByPodName(env, namespace, podName)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			common.GenFailureResult("metrics.getMetricsOfNamespaceByPodName.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}
