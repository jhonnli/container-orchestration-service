package harbor

import "github.com/gin-gonic/gin"

const HARBOR_NAME = "harbor"
const PROJECT_ID = "projectId"

func Init(engin *gin.Engine) []*gin.RouterGroup {
	data := make([]*gin.RouterGroup, 0)
	data = append(data, InitPorject(engin))
	data = append(data, InitRepository(engin))
	data = append(data, InitRepositoryOTag(engin))

	return data
}

func getHarborNameFromPath(ctx *gin.Context) string {
	return ctx.Param(HARBOR_NAME)
}

func getProjectIdFromPath(ctx *gin.Context) string {
	return ctx.Param(PROJECT_ID)
}
