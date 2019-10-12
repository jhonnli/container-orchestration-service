package proxy

import (
	"github.com/jhonnli/cmdb-service-sdk/businessEnvironment"
	"github.com/jhonnli/cmdb-service-sdk/client"
	"github.com/jhonnli/cmdb-service-sdk/cluster"
	harbor2 "github.com/jhonnli/cmdb-service-sdk/harbor"
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"github.com/jhonnli/container-orchestration-api/model/proxy"
	"github.com/jhonnli/container-orchestration-service/initial"
	"github.com/jhonnli/golib/logs"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var cmdbMutex sync.Mutex
var harborMutex sync.Mutex
var CmdbProxy *cmdbProxy

func Init() {
	cmdb := initial.Config.Cmdb
	CmdbProxy = NewCmdbProxy(cmdb.Address, cmdb.AppId, cmdb.Token, cmdb.IsDebug)
}

type cmdbProxy struct {
	clusterCache *cache.Cache
	harborCache  *cache.Cache
	client       *client.CmdbClient
}

func NewCmdbProxy(serverUrl string, appId int, appSecret string, isDebug bool) *cmdbProxy {
	client := client.NewCmdbClient(
		serverUrl,
		appId,
		appSecret,
		isDebug)

	return &cmdbProxy{
		clusterCache: cache.New(5*time.Minute, 6*time.Minute),
		harborCache:  cache.New(5*time.Minute, 6*time.Minute),
		client:       client,
	}
}

func (cmdb *cmdbProxy) GetK8sClusterInfo(env string) proxy.KubernetesAuthInfo {
	result, ok := cmdb.clusterCache.Get(env)
	if !ok {
		cmdbMutex.Lock()
		result, ok = cmdb.clusterCache.Get(env)
		if !ok {
			data, err := cmdb.getCluster(env)
			if err != nil {
				logs.Error("获取Cluster[ %s ] 发生异常, 在因: %s", env, err.Error())
			}
			result = proxy.KubernetesAuthInfo{
				Master: data.MasterUrl,
				Token:  data.AuthToken,
			}
			cmdb.clusterCache.SetDefault(env, result)
		}
		cmdbMutex.Unlock()
	}
	return result.(proxy.KubernetesAuthInfo)

}

func (cmdb *cmdbProxy) GetIngressInfo(cluster, group, project string) k8s.HPAParam {
	return k8s.HPAParam{}
}

func (cmdb *cmdbProxy) getEnv(env string) (businessEnvironment.BusinessEnvironmentInfo, error) {
	param := &businessEnvironment.GetBusinessEnvironmentInfoByEnvCodeRequest{
		EnvCode: env,
	}
	envResp, err := cmdb.client.BusinessEnvironment.GetBusinessEnvironmentInfoByEnvCode(param)
	if err != nil {
		logs.Error("从cmdb获取business environment [%s]失败，原因: %s", env, err.Error())
	}
	if envResp != nil {
		return envResp.Data, err
	} else {
		return businessEnvironment.BusinessEnvironmentInfo{}, err
	}
}

func (cmdb *cmdbProxy) getCluster(env string) (cluster.ClusterInfo, error) {
	envInfo, err := cmdb.getEnv(env)
	data, err := cmdb.client.Cluster.GetClusterInfoByEnvId(&cluster.GetClusterInfoByEnvIdRequest{
		EnvId: envInfo.EnvId,
	})
	if data != nil {
		return data.Data, err
	} else {
		return cluster.ClusterInfo{}, err
	}
}

func (cmdb cmdbProxy) getHarbor(env string) (proxy.HarborAuthInfo, error) {
	cluster, err := cmdb.getCluster(env)
	data, err := cmdb.client.Harbor.GetHarborInfo(&harbor2.GetHarborInfoRequest{
		HarborId: cluster.HarborId,
	})
	if err != nil {
		logs.Error("从cmdb获得%s下的harbor信息失败，原因: %s", env, err.Error())
	}
	pwdData, err := cmdb.client.Harbor.GetHarborLoginPwd(&harbor2.GetHarborLoginPwdRequest{
		HarborId: data.Data.HarborId,
	})
	if err != nil {
		logs.Error("从cmdb获得%s下的harbor密码失败，原因: %s", env, err.Error())
	}
	pwd := pwdData.Data
	harbor := data.Data
	return proxy.HarborAuthInfo{
		Server:   harbor.HarborUrl,
		Username: harbor.LoginUser,
		Password: pwd.LoginPwd,
	}, err
}

func (cmdb *cmdbProxy) GetHarborAuthInfo(env string) proxy.HarborAuthInfo {
	result, ok := cmdb.harborCache.Get(env)
	var err error
	if !ok {
		harborMutex.Lock()
		result, ok = cmdb.harborCache.Get(env)
		if !ok {
			result, err = cmdb.getHarbor(env)
			if err != nil {
				logs.Error("从cmdb中获取%s下harbor信息失败,原因: %s", env, err.Error())
			}
			cmdb.harborCache.SetDefault(env, result)
		}
		harborMutex.Unlock()
	}
	return result.(proxy.HarborAuthInfo)
}
