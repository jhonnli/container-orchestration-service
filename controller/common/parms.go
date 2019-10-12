package common

import "github.com/gin-gonic/gin"

const ENVNAME string = "env"

func GetEnvFromPath(ctx *gin.Context) string {
	return ctx.Param(ENVNAME)
}
