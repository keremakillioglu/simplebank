package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from config file or environment variable
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configurations from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app") // since app.env
	viper.SetConfigType("env") // json-xml etx could have been also be used

	viper.AutomaticEnv()
	// automatically override values from config file with the values of the env vars

	// start reading
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// unmarshal the values to target config object
	err = viper.Unmarshal(&config)
	return
}
