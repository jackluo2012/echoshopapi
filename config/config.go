package config

import "github.com/spf13/viper"

type weAppConfig struct {
	CodeToSessURL string
	AppID         string
	Secret        string
}
// WeAppConfig 微信小程序相关配置
var WeAppConfig weAppConfig

func readWxConfiguration() {
	viper.AddConfigPath("./server")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	WeAppConfig = &weAppConfig{
		CodeToSessURL:viper.GetString("weApp.CodeToSessURL"),
		AppID:viper.GetString("weApp.AppID"),
		Secret:viper.GetString("weApp.Secret"),
	}
}

func init() {
	readWxConfiguration()
}