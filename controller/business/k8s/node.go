package k8s

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s2 "github.com/jhonnli/container-orchestration-service/service/k8s"
	coreV1 "k8s.io/api/core/v1"
	"net/http"
)

var nodeService k8s.NodeInterface

func InitNodes() *gin.RouterGroup {
	nodeService = k8s2.NewNodeService()

	nodeApi := engine.Group("/v1/envs/:env/nodes")
	common.AddFilter(nodeApi)
	nodeApi.GET("/:nodeName", getNode)
	nodeApi.POST("/:nodeName/labels", addNodeLabels)
	nodeApi.DELETE("/:nodeName/labels", deleteNodeLabels)
	nodeApi.GET("", listNode)
	nodeApi.POST("/:nodeName/drains", drainPod)
	nodeApi.POST("/:nodeName/cordons", cordonNode)
	nodeApi.DELETE("/:nodeName/cordons", uncordonNode)
	return nodeApi
}

func drainPod(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "node.drain.param_error", env, nodeName) {
		return
	}
	err := nodeService.Drain(env, nodeName)

	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.drain.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func cordonNode(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "node.cordon.param_error", env, nodeName) {
		return
	}
	err := nodeService.Cordon(env, nodeName)

	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.drain.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func uncordonNode(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "node.uncordon.param_error", env, nodeName) {
		return
	}
	err := nodeService.Uncordon(env, nodeName)

	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.drain.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func getNode(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "node.get.param_error", env, nodeName) {
		return
	}
	data, err := nodeService.Get(env, nodeName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.get.error", err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.GenSuccessResult(data))
}

func listNode(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, "node.list.param_error", env) {
		return
	}
	//判断是否查询nodeLabels
	nodeLabels := context.DefaultQuery("nodeLabels", "")
	var data *coreV1.NodeList
	var err error
	if nodeLabels == "" {
		data, err = nodeService.List(env)
	} else {
		data, err = nodeService.ListNodeLabels(env, nodeLabels)
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.list.error", err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.GenSuccessResult(data))
}

func addNodeLabels(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "node.updateNodeLabels.param_error", env) {
		return
	}
	nodeLabels := make(map[string]string)
	err := common.GetJSONBody(context, &nodeLabels)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("node.updateNodeLabels.param_error", common.GetZhError(err)))
		return
	}
	err = nodeService.AddNodeLabels(env, nodeName, nodeLabels)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.updateNodeLabels.param_error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func deleteNodeLabels(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nodeName := getNodeNameFromPath(context)
	if common.ParamIsEmpty(context, "node.updateNodeLabels.param_error", env) {
		return
	}
	nodeLabels := make(map[string]string)
	err := common.GetJSONBody(context, &nodeLabels)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("node.updateNodeLabels.param_error", common.GetZhError(err)))
		return
	}
	err = nodeService.DeleteNodeLabels(env, nodeName, nodeLabels)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("node.updateNodeLabels.param_error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}
