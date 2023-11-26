package config

import "github.com/spf13/viper"

type SubscriberConfig struct {
	RpcUrl          string `mapstructure:"rpcUrl"`
	ContractAddress string `mapstructure:"contractAddress"`
	MongoUrl        string `mapstructure:"mongoUrl"`
}

func NewSubscriberConfig(filename string) (*SubscriberConfig, error) {
	if filename == "" {
		filename = "config.json"
	}

	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg SubscriberConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
