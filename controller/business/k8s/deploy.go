package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

var deployService k8s2.DeploymentInterface

func InitDeploy() *gin.RouterGroup {
	deployService = k8s3.NewDeploymentService()

	deployGroup := engine.Group("/v1/envs/:env/namespaces/:nsname/deployments")
	common.AddFilter(deployGroup)
	deployGroup.GET("/:deployName", getDeploy)
	deployGroup.PUT("", applyDeploy)
	deployGroup.POST("", createDeploy)
	deployGroup.PUT("/:deployName", updateDeploy)
	deployGroup.GET("", listDeploy)
	deployGroup.DELETE("/:deployName", deleteDeploy)
	return deployGroup
}

func getDeploy(contex *gin.Context) {
	nsname := getNamespaceNameFromPath(contex)
	deployName := getDeploymentNameFromPath(contex)
	env := common.GetEnvFromPath(contex)
	if common.ParamIsEmpty(contex, "deployment.get.param_error", env, deployName, nsname) {
		return
	}
	data, err := deployService.Get(env, nsname, deployName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("deployment.get.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func applyDeploy(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	if common.ParamIsEmpty(contex, "deployment.apply.param_error", env, nsname) {
		return
	}
	deployParam := &k8s.DeploymentParam{}
	err := common.GetJSONBody(contex, deployParam)
	if err != nil || nsname != deployParam.MetaData.Namespace || deployParam.MetaData.Name == "" {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("deployment.apply.param_error"))
		return
	}

	err = common.Validate.Struct(deployParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("deployment.apply.param_error", common.GetZhError(err)))
		return
	}

	data, err := deployService.Apply(env, *deployParam)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("deployment.create.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func createDeploy(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	if common.ParamIsEmpty(contex, "deployment.create.param_error", env, nsname) {
		return
	}
	deployParam := &k8s.DeploymentParam{}
	err := common.GetJSONBody(contex, deployParam)
	if err != nil || nsname != deployParam.MetaData.Namespace {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("deployment.create.param_error"))
		return
	}

	err = common.Validate.Struct(deployParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("deployment.create.param_error", common.GetZhError(err)))
		return
	}

	data, err := deployService.Create(env, *deployParam)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("deployment.create.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func updateDeploy(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	deployName := getDeploymentNameFromPath(contex)
	if common.ParamIsEmpty(contex, "deployment.update.param_error", env, nsname) {
		return
	}
	deployParam := &k8s.DeploymentParam{}
	err := common.GetJSONBody(contex, deployParam)
	if deployParam.MetaData.Name == "" {
		deployParam.MetaData.Name = deployName
	}
	if deployParam.MetaData.Namespace == "" {
		deployParam.MetaData.Namespace = nsname
	}
	if err != nil || nsname != deployParam.MetaData.Namespace || deployName != deployParam.MetaData.Name {
		contex.JSON(http.StatusOK, common.GenParamErrorResult("deployment.update.param_error"))
		return
	}

	err = common.Validate.Struct(deployParam)
	if err != nil {
		contex.JSON(http.StatusOK, common.GenFailureResult("deployment.update.param_error", common.GetZhError(err)))
		return
	}

	data, err := deployService.Update(env, *deployParam)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("deployment.update.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func listDeploy(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	if common.ParamIsEmpty(contex, "deployment.list.param_error", env, nsname) {
		return
	}
	data, err := deployService.List(env, nsname)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("deployment.list.error", err.Error()))
	} else {
		contex.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func deleteDeploy(contex *gin.Context) {
	env := common.GetEnvFromPath(contex)
	nsname := getNamespaceNameFromPath(contex)
	deployName := getDeploymentNameFromPath(contex)
	if common.ParamIsEmpty(contex, "deployment.delete.param_error", env, nsname) {
		return
	}
	err := deployService.Delete(env, nsname, deployName)
	if err != nil {
		contex.JSON(http.StatusInternalServerError, common.GenFailureResult("deployment.delete.error", err.Error()))
		return
	} else {
		contex.JSON(http.StatusOK, common.GenBoolResult())
	}
}
