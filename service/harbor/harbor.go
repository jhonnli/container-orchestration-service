package harbor

import (
	"encoding/base64"
	//"github.com/jhonnli/container-orchestration-api/model/proxy"
	"github.com/jhonnli/container-orchestration-service/service/common"
	//proxy2 "github.com/jhonnli/container-orchestration-service/service/proxy"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var harborMutex sync.Mutex

var HarborClient *harborClient

func Init() {
	HarborClient = NewHarborClient()
}

type harborClient struct {
}

func NewHarborClient() *harborClient {
	return &harborClient{}
}

//func (hs *harborClient) getHarborAuthInfo(harbor string) proxy.HarborAuthInfo {
//	return proxy2.CmdbProxy.GetHarborAuthInfo(harbor)
//}

func (hs harborClient) getURL(urlStr string) *url.URL {
	var scheme string
	var host string

	if strings.HasPrefix(urlStr, "https") {
		scheme = "https"
		host = string([]byte(urlStr)[8:])
	} else {
		scheme = "http"
		host = string([]byte(urlStr)[7:])
	}

	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

func (hs harborClient) GetClient(harbor string) *common.RestClient {
	//authInfo := hs.getHarborAuthInfo(harbor)
	/**
	@TODO 填充harbor的url信息，考虑从配置文件里面读取,填充用户名密码
	*/
	var harbor_url = " "
	client := common.NewRestClient(hs.getURL(harbor_url), 15)
	//client := common.NewRestClient(hs.getURL(authInfo.Server), 15)

	header := &http.Header{}
	var harbor_username = ""
	var harbor_password = ""
	header.Set("Authorization", "Basic "+basicAuth(harbor_username, harbor_password))
	//header.Set("Authorization", "Basic "+basicAuth(authInfo.Username, authInfo.Password))
	client.Header = header
	return client
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
