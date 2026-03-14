package config

import (
	"os"
	"sync"

	"github.com/spf13/viper"
)



var (
	once sync.Once
	configPath = "/home/phr/go-proj/mallive/internal/common/config"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}


func NewViperConfig() (err error) {
	once.Do(func ()  {
		viper.SetConfigName("global")
		viper.SetConfigType("yaml")

		// 支持通过环境变量指定配置路径
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			viper.AddConfigPath(envPath)
		}
		// 同时添加默认路径（本地开发用）
		viper.AddConfigPath(configPath)
		// 添加当前目录
		viper.AddConfigPath(".")
		// 添加 /config 目录（K8s 部署用）
		viper.AddConfigPath("/config")

		viper.AutomaticEnv()
		err = viper.ReadInConfig()
	})
	return
}