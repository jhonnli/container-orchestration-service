package harbor

import (
	harbor2 "github.com/jhonnli/container-orchestration-api/api/harbor"
	"github.com/jhonnli/container-orchestration-api/model/harbor"
)

func NewProjectService() harbor2.ProjectInterface {
	return &projectService{}
}

type projectService struct {
}

func (ps projectService) List(env string) ([]harbor.Project, error) {
	result := make([]harbor.Project, 0)
	err := HarborClient.GetClient(env).Get().RequestURI("/api/projects").
		Do().Into(&result)
	return result, err
}

func (ps projectService) Get(env, projectId string) (*harbor.Project, error) {
	result := &harbor.Project{}
	err := HarborClient.GetClient(env).Get().RequestURI("/api/projects/" + projectId).Do().Into(result)
	return result, err
}
