package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

var ingService k8s2.IngressInterface

func InitIngress() *gin.RouterGroup {
	ingService = k8s3.NewIngressService()

	ingressApi := engine.Group("/v1/envs/:env/namespaces/:nsname/ingresses")
	common.AddFilter(ingressApi)
	ingressApi.GET("/:ingName", getIngress)
	ingressApi.PUT("", applyIngress)
	ingressApi.POST("", createIngress)
	ingressApi.PUT("/:ingName", updateIngress)
	ingressApi.GET("", listIngress)
	ingressApi.DELETE("/:ingName", deleteIngress)

	return ingressApi
}

func getIngress(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	ingName := getIngressNameFromPath(contex)
	if common.ParamIsEmpty(contex, "ingress.get.param_error", env, nsname, ingName) {
		return
	}

	result, err := ingService.Get(env, nsname, ingName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("ingress.get.error", "获取Ingress异常"))
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(result))
	}
}

func applyIngress(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	if common.ParamIsEmpty(contex, "ingress.apply.param_error", env, nsname) {
		return
	}
	ingParam := &k8s.IngressParam{}
	err := common.GetJSONBody(contex, ingParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.apply.body_toobj_error", "参数解析错误"))
		return
	}

	if nsname != ingParam.MetaData.Namespace {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.apply.namespace_error", "参数校验失败"))
		return
	} else {
		ingParam.MetaData.Namespace = nsname
	}

	err = common.Validate.Struct(ingParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.apply.validate_error", common.GetZhError(err)))
		return
	}

	data, err := ingService.Apply(env, *ingParam)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("ingress.apply.error", err.Error()))
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func createIngress(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	if common.ParamIsEmpty(contex, "ingress.create.param_error", env, nsname) {
		return
	}
	ingParam := &k8s.IngressParam{}
	err := common.GetJSONBody(contex, ingParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.create.body_toobj_error", "参数解析错误"))
		return
	}
	if nsname != ingParam.MetaData.Namespace {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.create.namespace_error", "参数校验失败"))
		return
	} else {
		ingParam.MetaData.Namespace = nsname
	}

	err = common.Validate.Struct(ingParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.apply.validate_error", common.GetZhError(err)))
		return
	}

	data, err := ingService.Create(env, *ingParam)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("ingress.create.error", err.Error()))
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func updateIngress(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	ingName := getIngressNameFromPath(contex)
	if common.ParamIsEmpty(contex, "ingress.update.param_error", env, nsname, ingName) {
		return
	}
	ingParam := &k8s.IngressParam{}
	err := common.GetJSONBody(contex, ingParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.create.body_toobj_error", "参数解析错误"))
		return
	}
	if nsname != ingParam.MetaData.Namespace || ingName != ingParam.MetaData.Name {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.create.namespace_error", "参数校验失败"))
		return
	} else {
		ingParam.MetaData.Namespace = nsname
		ingParam.MetaData.Name = ingName
	}

	err = common.Validate.Struct(ingParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("ingress.apply.validate_error", common.GetZhError(err)))
		return
	}

	data, err := ingService.Update(env, *ingParam)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("ingress.update.error", err.Error()))
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func listIngress(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	if common.ParamIsEmpty(contex, "ingress.list.param_error", env, nsname) {
		return
	}
	data, err := ingService.List(env, nsname)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("ingress.list.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func deleteIngress(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	ingName := getIngressNameFromPath(contex)
	if common.ParamIsEmpty(contex, "ingress.delete.param_error", env, nsname, ingName) {
		return
	}
	err := ingService.Delete(env, nsname, ingName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("ingress.delete.error", "删除Ingress失败"))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenBoolResult())
	}
}
