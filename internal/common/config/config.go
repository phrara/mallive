package config

import (
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
		viper.AddConfigPath(configPath)
		viper.AutomaticEnv()
		err = viper.ReadInConfig()
	})
	return
}