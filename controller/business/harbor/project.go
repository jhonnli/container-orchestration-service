package harbor

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnli/container-orchestration-api/api/harbor"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	harbor2 "github.com/jhonnli/container-orchestration-service/service/harbor"
	"net/http"
)

var projectService harbor.ProjectInterface

func InitPorject(engin *gin.Engine) *gin.RouterGroup {
	projectService = harbor2.NewProjectService()

	projectApi := engin.Group("/v1/envs/:env/harbor/projects")
	common.AddFilter(projectApi)
	projectApi.GET("", listProject)
	projectApi.GET("/:projectId", getProject)
	return projectApi
}

func listProject(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, env) {
		return
	}
	data, err := projectService.List(env)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("project.list.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func getProject(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	projectId := getProjectIdFromPath(context)
	if common.ParamIsEmpty(context, env, projectId) {
		return
	}
	data, err := projectService.Get(env, projectId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("project.list.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}
