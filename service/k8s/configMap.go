package k8s

import (
	k8s2 "github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	apiV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

func NewConfigMapService() k8s2.ConfigMapInterface {
	return &configMapService{}
}

type configMapService struct {
}

func (cm *configMapService) getClient(env, namespace string) v1.ConfigMapInterface {
	return K8sClient.getClientSets(env).CoreV1().ConfigMaps(namespace)
}

func (cm configMapService) Create(env, namespace string, configMap k8s.ConfigMapParam) (*apiV1.ConfigMap, error) {
	requestData := cm.convertObject(configMap)
	return cm.getClient(env, namespace).Create(&requestData)
}

func (cm configMapService) List(env, namespace string) ([]apiV1.ConfigMap, error) {
	requestOptions := metaV1.ListOptions{}
	configList, err := cm.getClient(env, namespace).List(requestOptions)
	if err == nil {
		return configList.Items, err
	}
	return nil, err
}

func (cm configMapService) Update(env, namespace string, configMap k8s.ConfigMapParam) (*apiV1.ConfigMap, error) {
	requestData := cm.convertObject(configMap)
	return cm.getClient(env, namespace).Update(&requestData)
}

func (cm configMapService) convertObject(configMap k8s.ConfigMapParam) apiV1.ConfigMap {
	requestData := apiV1.ConfigMap{}
	requestData.Name = configMap.Name
	requestData.Data = configMap.Data
	return requestData
}

func (cm configMapService) Delete(env, namespace, configMapName string) error {
	return cm.getClient(env, namespace).Delete(configMapName, &metaV1.DeleteOptions{})
}

func (cm configMapService) Get(env, namespace, configMapName string) (*apiV1.ConfigMap, error) {
	return cm.getClient(env, namespace).Get(configMapName, metaV1.GetOptions{})
}
