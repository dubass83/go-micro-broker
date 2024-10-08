package util

import (
	"github.com/spf13/viper"
)

// Config store all configuration of the application
// the values read by viper from file or enviroment variables
type Config struct {
	Enviroment        string `mapstructure:"ENVIROMENT"`
	HTTPAddressString string `mapstructure:"HTTP_ADDRESS_STRING"`
	AuthService       string `mapstructure:"AUTH_SERVICE"`
}

// LoadConfig read configuration from file conf.env or enviroment variables
func LoadConfig(configPath string) (config Config, err error) {
	v := viper.New()
	v.SetConfigName("conf")
	v.SetConfigType("env")
	v.AddConfigPath(configPath)
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	v.AutomaticEnv()
	err = v.Unmarshal(&config)
	return
}