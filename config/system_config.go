package config

import "github.com/spf13/viper"

type SystemConfig struct {
	Host              string `mapstructure:"HOST"`
	Port              string `mapstructure:"PORT"`
	InitAdminUser     string `mapstructure:"INIT_ADMIN_USER"`
	InitAdminPassword string `mapstructure:"INIT_ADMIN_PASSWORD"`
}

func NewSystemConfig() (*SystemConfig, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	var cfg SystemConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
