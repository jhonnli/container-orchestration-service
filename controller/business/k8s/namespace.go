package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"github.com/jhonnli/golib/logs"
	"net/http"
)

var namespaceService k8s2.NamespaceInterface

func InitNamespace() *gin.RouterGroup {
	namespaceService = k8s3.NewNamespaceService()

	namespaceApi := engine.Group("/v1/envs/:env/namespaces")
	common.AddFilter(namespaceApi)
	namespaceApi.GET("/:nsname", getNamespace)
	namespaceApi.POST("", createNamespace)
	namespaceApi.GET("", listNamespace)
	namespaceApi.DELETE("/:nsname", deleteNamespace)

	return namespaceApi
}

func createNamespace(context *gin.Context) {
	env := common.GetEnvFromPath(context)

	if common.ParamIsEmpty(context, "namespace.create.param_error", env) {
		return
	}

	param := &k8s.NamespaceParam{}
	err := common.GetJSONBody(context, param)
	if err != nil {
		logs.Info("反序列化Namespace: %s 失败, 原因: %s\n", param.Name, err)
		context.JSON(http.StatusOK, common.GenParamErrorResult("namespace.create.body_parase"))
		return
	}

	err = common.Validate.Struct(param)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("namespace.create.body_parase", common.GetZhError(err)))
		return
	}
	resp, err := namespaceService.Create(env, *param)
	if err != nil {
		logs.Info("创建Namespace: %s 失败, 原因: %s\n", param.Name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("namespace.create.error", "创建Namespace失败"))
	} else {
		context.JSON(http.StatusOK, resp)
	}
}

func getNamespace(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	name := getNamespaceNameFromPath(context)
	if common.ParamIsEmpty(context, "namespace.get.param_error", env, name) {
		return
	}
	data, err := namespaceService.Get(env, name)
	if err != nil {
		logs.Info("获取Namespace: %s 失败, 原因: %s\n", name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("namespace.get", "获取Namespace失败"))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}

}

func listNamespace(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	data, err := namespaceService.List(env)
	if common.ParamIsEmpty(context, "namespace.list.param_error", env) {
		return
	}
	if err != nil {
		logs.Info("获取Namespace列表失败, 原因: %s\n", err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("namespace.list", "获取Namespace列表失败"))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func deleteNamespace(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	name := getNamespaceNameFromPath(context)

	if common.ParamIsEmpty(context, "namespace.delete.param_error", env, name) {
		return
	}
	err := namespaceService.Delete(env, name)
	if err != nil {
		logs.Info("删除Namespace: %s 失败, 原因: %s\n", name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("namespace.delete", "删除Namespace失败"))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(common.BoolResult{true}))
	}
}
