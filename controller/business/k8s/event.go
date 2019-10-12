package k8s

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	k8s2 "github.com/jhonnli/container-orchestration-service/service/k8s"
	"net/http"
)

const eventPathPrefix = "/v1/envs/:env/namespaces/:nsname"

var eventService k8s.EventInterface

func InitEvent() *gin.RouterGroup {
	eventService = k8s2.NewEventService()

	podOfNSApi := engine.Group(eventPathPrefix)
	common.AddFilter(podOfNSApi)
	podOfNSApi.GET("/events/:resourceName", getEvent)

	return podOfNSApi
}

func getEvent(context *gin.Context) {
	env := common.GetEnvFromPath(context)
	namespace := getNamespaceNameFromPath(context)
	resourceName := context.Param(RESOURCENAME)
	if common.ParamIsEmpty(context, env, namespace, resourceName) {
		return
	}
	data, err := eventService.Get(env, namespace, resourceName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.GenFailureResult("event.get.error", err.Error()))
	} else {
		context.JSON(http.StatusOK, common.GenSuccessResult(data))
	}
}
