package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBSource       string `mapstructure:"DB_SOURCE"`
	SevrverAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path, filename string) (config Config, err error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return config, err
}
