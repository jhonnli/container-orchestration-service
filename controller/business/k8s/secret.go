package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

var secretService k8s2.SecretInterface

func InitSecret() *gin.RouterGroup {
	secretService = k8s3.NewSecretService()

	secretApi := engine.Group("/v1/envs/:env/namespaces/:nsname/secrets")
	common.AddFilter(secretApi)
	secretApi.GET("", listSecret)
	secretApi.GET("/:secretName", getSecret)
	secretApi.POST("", createSecret)
	secretApi.PUT("/:secretName", updateSecret)
	secretApi.DELETE("/:secretName", deleteSecret)
	return secretApi
}

func listSecret(context *gin.Context) {
	nsname := getNamespaceNameFromPath(context)
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, "secret.list.param_error", env, nsname) {
		return
	}
	data, err := secretService.List(env, nsname)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("secret.list.error", err.Error()))
		return
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func createSecret(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	if common.ParamIsEmpty(contex, "secret.create.param_error", env, nsname) {
		return
	}
	secretParm := &k8s.SecretParam{}
	err := common.GetJSONBody(contex, secretParm)
	if err != nil || secretParm.Name == "" || secretParm.Data == nil || len(secretParm.Data) == 0 {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("secret.apply.param_error"))
		return
	}
	data, err := secretService.Create(env, nsname, *secretParm)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("secret.create.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func updateSecret(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	secretName := getSecretNameFromPath(contex)
	if common.ParamIsEmpty(contex, "secret.update.param_error", env, nsname, secretName) {
		return
	}
	secretParm := &k8s.SecretParam{}
	err := common.GetJSONBody(contex, secretParm)
	if err != nil || secretParm.Name != secretName || secretParm.Data == nil || len(secretParm.Data) == 0 {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("secret.apply.param_error"))
		return
	}
	data, err := secretService.Update(env, nsname, *secretParm)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("secret.update.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func deleteSecret(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	secretName := getSecretNameFromPath(contex)
	if common.ParamIsEmpty(contex, "secret.update.param_error", env, nsname, secretName) {
		return
	}
	err := secretService.Delete(env, nsname, secretName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("secret.delete.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenBoolResult())
	}
}

func getSecret(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	secretName := getSecretNameFromPath(contex)
	if common.ParamIsEmpty(contex, "secret.update.param_error", env, nsname, secretName) {
		return
	}
	data, err := secretService.Get(env, nsname, secretName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("secret.delete.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}
