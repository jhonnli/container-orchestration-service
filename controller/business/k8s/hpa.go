package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

var hpaService k8s2.HPAInterface

func InitHPA() *gin.RouterGroup {
	hpaService = k8s3.NewHPAService()

	hpaApi := engine.Group("/v1/envs/:env/namespaces/:nsname/hpas")
	common.AddFilter(hpaApi)
	hpaApi.GET("/:hpaName", getHPA)
	hpaApi.PUT("", applyHPA)
	hpaApi.POST("", createHPA)
	hpaApi.PUT("/:hpaName", updateHPA)
	hpaApi.GET("", listHPA)
	hpaApi.DELETE("/:hpaName", deleteHPA)

	return hpaApi
}

func applyHPA(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	hpaParam := &k8s.HPAParam{}
	err := common.GetJSONBody(context, hpaParam)
	if common.ParamIsEmpty(context, "hpa.apply.param_error", env, nsname) {
		return
	}
	if nsname != hpaParam.MetaData.Namespace || hpaParam.MetaData.Name == "" {
		context.JSON(http.StatusOK, common.GenParamErrorResult("hpa.apply.param_error"))
		return
	}

	err = common.Validate.Struct(hpaParam)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("hpa.apply.param_error", common.GetZhError(err)))
		return
	}

	data, err := hpaService.Apply(env, *hpaParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("hpa.apply.error", err.Error()))
		return
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func getHPA(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	hpaName := getHpaNameFromPath(context)
	if common.ParamIsEmpty(context, "hpa.get.param_error", env, nsname) {
		return
	}
	data, err := hpaService.Get(env, nsname, hpaName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("hpa.get.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func createHPA(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	if common.ParamIsEmpty(context, "hpa.create.param_error", env, nsname) {
		return
	}
	hpaParam := &k8s.HPAParam{}
	err := common.GetJSONBody(context, hpaParam)

	if err != nil || nsname != hpaParam.MetaData.Namespace {
		context.JSON(http.StatusOK, common.GenParamErrorResult("hpa.create.param_error"))
		return
	}

	err = common.Validate.Struct(hpaParam)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("hpa.create.param_error", common.GetZhError(err)))
		return
	}

	data, err := hpaService.Create(env, *hpaParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("hpa.create.error", err.Error()))
		return
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}

}

func listHPA(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	if common.ParamIsEmpty(context, "hpa.list.param_error", env, nsname) {
		return
	}
	data, err := hpaService.List(env, nsname)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("hpa.list.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func updateHPA(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	hpaName := getHpaNameFromPath(context)
	if common.ParamIsEmpty(context, "hpa.update.param_error", env, nsname) {
		return
	}
	hpaParam := &k8s.HPAParam{}
	err := common.GetJSONBody(context, hpaParam)
	if err != nil || nsname != hpaParam.MetaData.Namespace {
		context.JSON(http.StatusOK, common.GenParamErrorResult("hpa.update.param_error"))
		return
	}
	if hpaParam.MetaData.Name == "" {
		hpaParam.MetaData.Name = hpaName
	}

	err = common.Validate.Struct(hpaParam)
	if err != nil {
		context.JSON(http.StatusOK, common.GenFailureResult("hpa.update.param_error", common.GetZhError(err)))
		return
	}

	data, err := hpaService.Update(env, *hpaParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("hpa.update.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}

}

func deleteHPA(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	nsname := getNamespaceNameFromPath(context)
	hpaName := getHpaNameFromPath(context)
	if common.ParamIsEmpty(context, "hpa.delete.param_error", env, nsname, hpaName) {
		return
	}
	err := hpaService.Delete(env, nsname, hpaName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("hpa.delete.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenBoolResult())
	}
}
