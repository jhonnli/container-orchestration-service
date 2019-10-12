module container-orchestration-service

go 1.13

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-gonic/gin v1.4.0
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/jhonnli/container-orchestration-api v0.0.0-20191011095017-7320cd233f58
	github.com/jhonnli/logs v0.0.0-20191011085002-405f78123657
	github.com/json-iterator/go v1.1.7
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/spf13/viper v1.4.0
	gopkg.in/go-playground/validator.v9 v9.30.0
	k8s.io/api v0.0.0-20191010143144-fbf594f18f80
	k8s.io/apimachinery v0.0.0-20191006235458-f9f2f3f8ab02
	k8s.io/client-go v11.0.0+incompatible
)
