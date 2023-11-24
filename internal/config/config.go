package config

import "github.com/spf13/viper"

type AuthServiceConfig struct {
	IPFSUrl         string `mapstructure:"ipfsUrl"`
	RPCUrl          string `mapstructure:"rpcUrl"`
	ContractAddress string `mapstructure:"contractAddress"`
	ResolverPrefix  string `mapstructure:"resolverPrefix"`
	KeyDir          string `mapstructure:"keyDir"`
	RedisUrl        string `mapstructure:"redisUrl"`
}

type DidServiceConfig struct {
	RPCUrl    string `mapstructure:"rpcUrl"`
	IssuerUrl string `mapstructure:"issuerUrl"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	RedisUrl  string `mapstructure:"redisUrl"`
}

type IpfsServiceConfig struct {
	ApiKey string `mapstructure:"apiKey"`
	Secret string `mapstructure:"secret"`
}

type JwtServiceConfig struct {
	Issuer      string `mapstructure:"issuer"`
	Audience    string `mapstructure:"audience"`
	SecretKey   string `mapstructure:"secretKey"`
	TokenExpiry int64  `mapstructure:"tokenExpiry"`
}

type NoticeServiceConfig struct {
	DbName     string `mapstructure:"dbName"`
	Collection string `mapstructure:"collection"`
}

type TicketServiceConfig struct {
}

type UserServiceConfig struct {
	DbName     string `mapstructure:"dbName"`
	Collection string `mapstructure:"collection"`
}

type Config struct {
	Auth      AuthServiceConfig   `mapstructure:"auth"`
	Did       DidServiceConfig    `mapstructure:"did"`
	Ipfs      IpfsServiceConfig   `mapstructure:"ipfs"`
	Jwt       JwtServiceConfig    `mapstructure:"jwt"`
	Notice    NoticeServiceConfig `mapstructure:"notice"`
	Ticket    TicketServiceConfig `mapstructure:"ticket"`
	User      UserServiceConfig   `mapstructure:"user"`
	MongoUrl  string              `mapstructure:"mongoUrl"`
	ServerUrl string              `mapstructure:"serverUrl"`
}

func NewConfig(filename string) (*Config, error) {
	if filename == "" {
		filename = "config.json"
	}

	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
