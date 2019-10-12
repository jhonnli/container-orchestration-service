package k8s

import (
	"github.com/gin-gonic/gin"
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s3 "github.com/jhonnli/container-orchestration-service/service/k8s"
	//"github.com/jhonnli/golib/logs"
	"github.com/jhonnli/logs"
	"net/http"
)

var applicationService k8s2.ApplicationInterface

func InitApplication() {
	applicationService = k8s3.NewApplicationService()

	idxApi := engine.Group("/api/v1/envs/:env/applications")
	common.AddApiFilter(idxApi)
	idxApi.PUT("", applyApplication)
	idxApi.DELETE("/:appName", deleteApplication)
}

func deleteApplication(context *gin.Context) {

}

func applyApplication(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	param := &k8s.ApplicationParam{}
	err := common.GetJSONBody(context, param)
	if err != nil {
		logs.Info(err)
		context.JSON(http.StatusOK, common.GenFailureResult("application.apply.body_parse_error", err.Error()))
		return
	}
	err = common.Validate.Struct(param)
	if err != nil {
		logs.Info(err)
		context.JSON(http.StatusOK, common.GenFailureResult("application.apply.body_parse_error", common.GetZhError(err)))
		return
	}
	result := applicationService.Apply(env, *param)
	context.JSON(http.StatusOK, common.GenSuccessResult(result))
}
