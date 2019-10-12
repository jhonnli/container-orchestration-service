package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"github.com/jhonnli/logs"
	"net/http"
)

var svcService k8s2.ServiceInterface

func InitService() *gin.RouterGroup {
	svcService = k8s3.NewSvcService()

	svcApi := engine.Group("/v1/envs/:env/namespaces/:nsname/services")
	common.AddFilter(svcApi)
	svcApi.GET("/:svcname", getService)
	svcApi.PUT("", applyService)
	svcApi.POST("", createService)
	svcApi.PUT("/:svcname", updateService)
	svcApi.GET("", listService)
	svcApi.DELETE("/:svcname", deleteService)

	return svcApi
}

func createService(context *gin.Context) {
	clusterName := common.GetEnvFromPath(context)
	namespace := common.GetStringParam(context, "nsname")
	if common.ParamIsEmpty(context, "service.create.param_error", clusterName, namespace) {
		return
	}
	param := &k8s.ServiceParam{}
	common.GetJSONBody(context, param)
	err := common.Validate.Struct(param)

	if err != nil {
		context.JSON(http.StatusBadRequest, common.GenFailureResult("service.create.param_error", common.GetZhError(err)))
		return
	}
	if param.MetaData.Namespace != namespace {
		context.JSON(http.StatusOK, common.GenFailureResult("service.apply.param.error", "保存service参数错误"))
		return
	}

	resp, err := svcService.Create(clusterName, *param)
	if err != nil {
		logs.Info("创建%s.Service: %s 失败, 原因: %s\n", namespace, param.MetaData.Name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("Service.create.error", "创建Service失败"))
	} else {
		context.JSON(http.StatusOK, resp)
	}
}

func updateService(context *gin.Context) {
	clusterName := common.GetEnvFromPath(context)
	namespace := common.GetStringParam(context, "nsname")
	name := common.GetStringParam(context, "svcname")
	if common.ParamIsEmpty(context, "service.update.param_error", clusterName, namespace, name) {
		return
	}
	param := &k8s.ServiceParam{}
	common.GetJSONBody(context, param)
	err := common.Validate.Struct(param)
	if err != nil {
		context.JSON(http.StatusBadRequest, common.GenFailureResult("service.update.param_error", common.GetZhError(err)))
		return
	}

	if param.MetaData.Name != name || param.MetaData.Namespace != "" || param.MetaData.Namespace != namespace {
		context.JSON(http.StatusOK, common.GenFailureResult("service.apply.param.error", "保存service参数错误"))
		return
	}

	param.MetaData.Namespace = namespace
	param.MetaData.Name = name

	resp, err := svcService.Create(clusterName, *param)
	if err != nil {
		logs.Info("创建%s.Service: %s 失败, 原因: %s\n", namespace, param.MetaData.Name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("Service.create.error", "创建Service失败"))
	} else {
		context.JSON(http.StatusOK, resp)
	}
}

func applyService(context *gin.Context) {
	clusterName := common.GetEnvFromPath(context)
	namespace := common.GetStringParam(context, "nsname")

	if common.ParamIsEmpty(context, "service.apply.param_error", clusterName, namespace) {
		return
	}

	param := &k8s.ServiceParam{}
	err := common.GetJSONBody(context, param)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("service.apply.param.error", "获取json body失败"))
		return
	}
	err = common.Validate.Struct(param)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("service.apply.param.error", common.GetZhError(err)))
		return
	}

	if param.MetaData.Namespace != namespace {
		context.JSON(http.StatusOK, common.GenFailureResult("service.apply.param.error", "保存service参数错误"))
		return
	}

	resp, err := svcService.Apply(clusterName, *param)
	if err != nil {
		logs.Info("创建%s.Service: %s 失败, 原因: %s\n", namespace, param.MetaData.Name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("Service.create.error", "创建Service失败"))
	} else {
		context.JSON(http.StatusOK, resp)
	}
}

func getService(context *gin.Context) {
	clusterName := common.GetEnvFromPath(context)
	namespace := common.GetStringParam(context, "nsname")
	name := common.GetStringParam(context, "svcname")
	if common.ParamIsEmpty(context, "service.get.param_error", clusterName, namespace, name) {
		return
	}
	data, err := svcService.Get(clusterName, namespace, name)
	if err != nil {
		logs.Info("获取Service: %s 失败, 原因: %s\n", name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("Service.get", "获取Service失败"))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}

}

func listService(context *gin.Context) {
	clusterName := common.GetEnvFromPath(context)
	namespace := common.GetStringParam(context, "nsname")
	if common.ParamIsEmpty(context, "service.list.param_error", clusterName, namespace) {
		return
	}
	data, err := svcService.List(clusterName, namespace)
	if err != nil {
		logs.Info("获取Service列表失败, 原因: %s\n", err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("Service.list", "获取Service列表失败"))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func deleteService(context *gin.Context) {
	clusterName := common.GetEnvFromPath(context)
	namespace := common.GetStringParam(context, "nsname")
	name := common.GetStringParam(context, "svcname")
	if common.ParamIsEmpty(context, "service.delete.param_error", clusterName, namespace) {
		return
	}
	err := svcService.Delete(clusterName, namespace, name)
	if err != nil {
		logs.Info("删除Service: %s 失败, 原因: %s\n", name, err)
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("Service.delete", "删除Service失败"))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(common.BoolResult{true}))
	}
}
