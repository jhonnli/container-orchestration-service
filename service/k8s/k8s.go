package k8s

import (
	"fmt"
	//"github.com/jhonnli/container-orchestration-service/service/proxy"
	"github.com/jhonnli/logs"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

var k8sMutex sync.Mutex

var K8sClient *k8sClient

func Init() {
	K8sClient = NewK8sClient()
}

type k8sClient struct {
	clientsetMap map[string]*kubernetes.Clientset
}

type RestClientParam struct {
	pathPrefix string
	group      string
	version    string
}

func (kc *k8sClient) getClientset(env string) *kubernetes.Clientset {
	var client *kubernetes.Clientset
	var ok bool
	var err error
	client, ok = kc.clientsetMap[env]
	if !ok {
		k8sMutex.Lock()
		client, ok = kc.clientsetMap[env]
		if !ok {
			client, err = kubernetes.NewForConfig(kc.getConfig(env))
			if err != nil {
				fmt.Println("生成clientset失败:", err)
			}
			kc.clientsetMap[env] = client
		}

		k8sMutex.Unlock()
	}
	return client
}

func (kc *k8sClient) getConfig(env string) *rest.Config {
	//@TODO 原从cmdb中读取k8s的配置信息需要改成从配置文件中读取
	//authInfo := proxy.CmdbProxy.GetK8sClusterInfo(env)
	var k8s_master = ""
	var k8s_token = ""
	config, err := clientcmd.BuildConfigFromFlags(k8s_master, "")
	if err != nil {
		panic(err)
	}
	config.Insecure = true
	config.BearerToken = k8s_token
	return config
}

func (kc *k8sClient) getClientSets(clusterName string) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(kc.getConfig(clusterName))
	if err != nil {
		fmt.Println("生成clientset失败:", err)
	}
	return clientset
}

func (kc *k8sClient) getRestClient(env string, param RestClientParam) *rest.RESTClient {
	config := kc.getConfig(env)
	config.GroupVersion = &schema.GroupVersion{
		Group:   param.group,
		Version: param.version,
	}
	config.APIPath = param.pathPrefix
	if config.APIPath == "" {
		config.APIPath = "/apis"
	}
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: scheme.Codecs,
	}
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return restClient
}

func NewK8sClient() *k8sClient {
	client := &k8sClient{
		clientsetMap: make(map[string]*kubernetes.Clientset),
	}
	logs.Info("初始化k8s client success")
	return client
}
