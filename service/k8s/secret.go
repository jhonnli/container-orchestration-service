package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	apiV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

func NewSecretService() k8s2.SecretInterface {
	return &secretService{}
}

type secretService struct {
}

func (s *secretService) getClient(env, namespace string) v1.SecretInterface {
	return K8sClient.getClientSets(env).CoreV1().Secrets(namespace)
}

func (s secretService) Create(env, namespace string, secretParam k8s.SecretParam) (*apiV1.Secret, error) {
	requestData := s.convertObject(secretParam)
	return s.getClient(env, namespace).Create(&requestData)
}

func (s secretService) List(env, namespace string) ([]apiV1.Secret, error) {
	requestOptions := metaV1.ListOptions{}
	configList, err := s.getClient(env, namespace).List(requestOptions)
	if err == nil {
		return configList.Items, err
	}
	return nil, err
}

func (s secretService) Update(env, namespace string, secret k8s.SecretParam) (*apiV1.Secret, error) {
	requestData := s.convertObject(secret)
	return s.getClient(env, namespace).Update(&requestData)
}

func (s secretService) convertObject(secretParam k8s.SecretParam) apiV1.Secret {
	stringMap := secretParam.Data
	byteMap := make(map[string][]byte)
	for k, v := range stringMap {
		byteMap[k] = []byte(v)
	}
	requestData := apiV1.Secret{}
	requestData.Name = secretParam.Name
	requestData.Data = byteMap
	return requestData
}

func (s secretService) Delete(env, namespace, secretName string) error {
	return s.getClient(env, namespace).Delete(secretName, &metaV1.DeleteOptions{})
}

func (s secretService) Get(env, namespace, secretName string) (*apiV1.Secret, error) {
	return s.getClient(env, namespace).Get(secretName, metaV1.GetOptions{})
}
