package config

import (
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init takes the environment, starts the viper
// and loads the corresponding configuration file.
func Init(env string, dir string) (err error) {

	config = viper.New()

	config.SetConfigType("json")
	config.SetConfigName(env)
	config.AddConfigPath(dir + "/config/")

	err = config.ReadInConfig()
	if err != nil {
		return err
	}

	return
}

//Config returns the configuration struct
func Config() *viper.Viper {

	return config
}
