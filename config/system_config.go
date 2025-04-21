package config

import "github.com/spf13/viper"

type SystemConfig struct {
	SysAdminUser     string `mapstructure:"SYSTEM_ADMIN_USER"`
	SysAdminPassword string `mapstructure:"SYSTEM_ADMIN_PASSWORD"`
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
