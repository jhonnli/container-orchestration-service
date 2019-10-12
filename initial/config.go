package initial

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jhonnli/logs"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Config config

func InitConfig() {
	DEPLOY_ENV := os.Getenv("DEPLOY_ENV")
	switch DEPLOY_ENV {
	case "dev":
		viper.SetConfigFile("config_dev.yaml")
	case "test":
		viper.SetConfigFile("config_test.yaml")
	case "prod":
		viper.SetConfigFile("config_prod.yaml")
	default:
		log.Fatalln("没有找到环境变量DEPLOY_ENV")
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		logs.Error(err)
	}

	conf, err := parseConfig()
	if err != nil {
		os.Exit(1)
	}

	Config = *conf

	viper.WatchConfig()
	viper.OnConfigChange(func(event fsnotify.Event) {
		if event.Op == fsnotify.Write {
			conf, err := parseConfig()
			if err != nil {
				return
			}
			Config = *conf
			log.Println("config reload.")
		}

	})
}

func parseConfig() (*config, error) {
	var conf config
	err := viper.Unmarshal(&conf)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return &conf, nil
}

type config struct {
	Listen listen
	Cmdb   cmdb
}

type cmdb struct {
	Address string
	AppId   int
	Token   string
	IsDebug bool
}

type listen struct {
	Domain  string `mapstructure:"domain"`
	Address string `mapstructure:"address"`
}
