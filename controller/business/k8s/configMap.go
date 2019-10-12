package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

var configMapService k8s2.ConfigMapInterface

func InitConfigMap() *gin.RouterGroup {
	configMapService = k8s3.NewConfigMapService()
	configMapApi := engine.Group("/v1/envs/:env/namespaces/:nsname/configMaps")
	common.AddFilter(configMapApi)
	configMapApi.GET("", listConfigMap)
	configMapApi.GET("/:configMapName", getConfigMap)
	configMapApi.POST("", createConfigMap)
	configMapApi.PUT("/:configMapName", updateConfigMap)
	configMapApi.DELETE("/:configMapName", deleteConfigMap)
	return configMapApi
}

func listConfigMap(context *gin.Context) {
	nsname := getNamespaceNameFromPath(context)
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, "configMap.list.param_error", env, nsname) {
		return
	}
	data, err := configMapService.List(env, nsname)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("configMap.list.error", err.Error()))
		return
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func createConfigMap(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	if common.ParamIsEmpty(contex, "configMap.create.param_error", env, nsname) {
		return
	}
	configMapParm := &k8s.ConfigMapParam{}
	err := common.GetJSONBody(contex, configMapParm)
	if err != nil || configMapParm.Name == "" || configMapParm.Data == nil || len(configMapParm.Data) == 0 {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("configMap.apply.param_error"))
		return
	}
	data, err := configMapService.Create(env, nsname, *configMapParm)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("configMap.create.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func updateConfigMap(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	configMapName := getConfigMapNameFromPath(contex)
	if common.ParamIsEmpty(contex, "configMap.update.param_error", env, nsname, configMapName) {
		return
	}
	configMapParm := &k8s.ConfigMapParam{}
	err := common.GetJSONBody(contex, configMapParm)
	if err != nil || configMapParm.Name != configMapName || configMapParm.Data == nil || len(configMapParm.Data) == 0 {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("configMap.apply.param_error"))
		return
	}
	data, err := configMapService.Update(env, nsname, *configMapParm)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("configMap.update.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func deleteConfigMap(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	configMapName := getConfigMapNameFromPath(contex)
	if common.ParamIsEmpty(contex, "configMap.delete.param_error", env, nsname, configMapName) {
		return
	}
	err := configMapService.Delete(env, nsname, configMapName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("configMap.delete.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func getConfigMap(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	configMapName := getConfigMapNameFromPath(contex)
	if common.ParamIsEmpty(contex, "configMap.delete.param_error", env, nsname, configMapName) {
		return
	}
	data, err := configMapService.Get(env, nsname, configMapName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("configMap.delete.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}
