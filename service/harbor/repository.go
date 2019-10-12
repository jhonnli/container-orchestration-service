package harbor

import (
	harbor2 "github.com/jhonnli/container-orchestration-api/api/harbor"
	"github.com/jhonnli/container-orchestration-api/model/harbor"
	"github.com/jhonnli/logs"
)

func NewRepositoryService() harbor2.RepositoryInterface {
	return &repositoryService{
		projectService: NewProjectService(),
	}
}

type repositoryService struct {
	projectService harbor2.ProjectInterface
}

func (rs repositoryService) List(env string, projectId string) ([]harbor.Repository, error) {
	result := make([]harbor.Repository, 0)
	err := HarborClient.GetClient(env).Get().RequestURI("/api/repositories").
		Param("project_id", projectId).Do().Into(&result)
	return result, err
}

func (rs repositoryService) ListTag(env, projectId, repoName string) ([]harbor.DetailedTag, error) {
	project, err := rs.projectService.Get(env, projectId)
	if err != nil {
		logs.Info("获得%s环境下%s项目下%s仓库失败，原因:", env, projectId, repoName, err.Error())
		return nil, err
	}
	repoFullName := project.Name + "/" + repoName
	result := make([]harbor.DetailedTag, 0)
	err = HarborClient.GetClient(env).Get().RequestURI("/api/repositories/" + repoFullName + "/tags").Do().Into(&result)
	return result, err
}
