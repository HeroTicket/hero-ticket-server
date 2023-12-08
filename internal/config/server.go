package config

import "github.com/spf13/viper"

type AuthServiceConfig struct {
	IPFSUrl         string `mapstructure:"ipfsUrl"`
	ContractAddress string `mapstructure:"contractAddress"`
	ResolverPrefix  string `mapstructure:"resolverPrefix"`
	KeyDir          string `mapstructure:"keyDir"`
	RedisUrl        string `mapstructure:"redisUrl"`
}

type DidServiceConfig struct {
	IssuerUrl string `mapstructure:"issuerUrl"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	RedisUrl  string `mapstructure:"redisUrl"`
	DbName    string `mapstructure:"dbName"`
}

type IpfsServiceConfig struct {
	ApiKey string `mapstructure:"apiKey"`
	Secret string `mapstructure:"secret"`
}

type JwtServiceConfig struct {
	Issuer             string `mapstructure:"issuer"`
	Audience           string `mapstructure:"audience"`
	AccessTokenKey     string `mapstructure:"accessTokenKey"`
	RefreshTokenKey    string `mapstructure:"refreshTokenKey"`
	AccessTokenExpiry  int64  `mapstructure:"accessTokenExpiry"`
	RefreshTokenExpiry int64  `mapstructure:"refreshTokenExpiry"`
}

type NoticeServiceConfig struct {
	DbName string `mapstructure:"dbName"`
}

type TicketServiceConfig struct {
	DbName          string `mapstructure:"dbName"`
	ContractAddress string `mapstructure:"contractAddress"`
	PrivateKey      string `mapstructure:"privateKey"`
	MoralisApiKey   string `mapstructure:"moralisApiKey"`
}

type UserServiceConfig struct {
	DbName string `mapstructure:"dbName"`
}

type ServerConfig struct {
	Auth      AuthServiceConfig   `mapstructure:"auth"`
	Did       DidServiceConfig    `mapstructure:"did"`
	Ipfs      IpfsServiceConfig   `mapstructure:"ipfs"`
	Jwt       JwtServiceConfig    `mapstructure:"jwt"`
	Notice    NoticeServiceConfig `mapstructure:"notice"`
	Ticket    TicketServiceConfig `mapstructure:"ticket"`
	User      UserServiceConfig   `mapstructure:"user"`
	MongoUrl  string              `mapstructure:"mongoUrl"`
	RpcUrl    string              `mapstructure:"rpcUrl"`
	ServerUrl string              `mapstructure:"serverUrl"`
}

func NewServerConfig(filename string) (*ServerConfig, error) {
	if filename == "" {
		filename = "config.json"
	}

	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg ServerConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
