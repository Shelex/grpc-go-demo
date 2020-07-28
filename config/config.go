package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Host              string
		ClientAPIPort     string
		DomainServicePort string
		MongoURL          string
		PathToTLSCertFile string
		PathToTLSKeyFile  string
		Env               string
		LogLevel          string
		PrettyLogOutput   bool
		InitTimeout       time.Duration
	}
)

func GetEnv() Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("CLIENT_API_PORT", "8080")
	viper.SetDefault("DOMAIN_SERVICE_PORT", "9000")
	viper.SetDefault("MONGO_URL", "localhost:27017")
	viper.SetDefault("PATH_TO_TLS_CERT", "cert.pem")
	viper.SetDefault("PATH_TO_TLS_KEY", "key.pem")
	viper.SetDefault("ENV", "production")
	viper.SetDefault("PRETTY_LOG_OUTPUT", true)
	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("INIT_TIMEOUT", 10*time.Second)

	return Config{
		Host:              viper.GetString("HOST"),
		ClientAPIPort:     viper.GetString("CLIENT_API_PORT"),
		DomainServicePort: viper.GetString("DOMAIN_SERVICE_PORT"),
		MongoURL:          viper.GetString("MONGO_URL"),
		PathToTLSCertFile: viper.GetString("PATH_TO_TLS_CERT"),
		PathToTLSKeyFile:  viper.GetString("PATH_TO_TLS_KEY"),
		Env:               viper.GetString("ENV"),
		LogLevel:          viper.GetString("LOG_LEVEL"),
		PrettyLogOutput:   viper.GetBool("PRETTY_LOG_OUTPUT"),
		InitTimeout:       viper.GetDuration("INIT_TIMEOUT"),
	}
}
