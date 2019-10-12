package controller

import (
	"github.com/gin-gonic/gin"
	apik8s "github.com/jhonnli/container-orchestration-service/controller/api/k8s"
	harborController "github.com/jhonnli/container-orchestration-service/controller/business/harbor"
	bk8s "github.com/jhonnli/container-orchestration-service/controller/business/k8s"
	"github.com/jhonnli/container-orchestration-service/controller/common"
	"net/http"
)

func Init(engine *gin.Engine) {
	engine.GET("/", index)
	engine.Any("/health", health)
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, common.Result{Code: "not found", Message: "Page not found"})
	})
	bk8s.Init(engine)
	harborController.Init(engine)
	apik8s.Init(engine)
}

func index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "I'm OK")
}

func health(ctx *gin.Context) {
	result := make(map[string]string)
	result["status"] = "UP"
	ctx.JSON(http.StatusOK, result)
}
