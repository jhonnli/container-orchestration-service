package harbor

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnli/container-orchestration-api/api/harbor"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	harbor2 "github.com/jhonnli/container-orchestration-service/service/harbor"
	"net/http"
)

var repositoryService harbor.RepositoryInterface

func InitRepository(engin *gin.Engine) *gin.RouterGroup {
	repositoryService = harbor2.NewRepositoryService()

	repoApi := engin.Group("/v1/envs/:env/harbor/projects/:projectId/repositories")
	common.AddFilter(repoApi)
	repoApi.GET("", listRepository)

	return repoApi
}
func InitRepositoryOTag(engin *gin.Engine) *gin.RouterGroup {
	tagApi := engin.Group("/v1/envs/:env/harbor/projects/:projectId/repositories/:repoName/tags")
	tagApi.GET("", listRepositoryTag)

	return tagApi
}

func listRepository(context *gin.Context) {
	projectId := getProjectIdFromPath(context)
	env := common.GetEnvFromPath(context)
	if common.ParamIsEmpty(context, env, projectId) {
		return
	}
	data, err := repositoryService.List(env, projectId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("repository.list.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}

func listRepositoryTag(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	projectId := getProjectIdFromPath(context)
	repoName := context.Param("repoName")
	if common.ParamIsEmpty(context, env, projectId, repoName) {
		return
	}
	data, err := repositoryService.ListTag(env, projectId, repoName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("repository.tag.list.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}
