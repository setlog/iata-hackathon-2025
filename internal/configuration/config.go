package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	GcProjectId       string   `mapstructure:"GCLOUD_PROJECT_ID"`
	GcLocation        string   `mapstructure:"GCLOUD_LOCATION"`
	GcBucketName      string   `mapstructure:"GCLOUD_BUCKETNAME"`
	AiModel           string   `mapstructure:"AI_MODEL"`
	OAuthIssuer       string   `mapstructure:"OAUTH_ISSUER"`
	OAuthClientIds    []string `mapstructure:"OAUTH_CLIENT_IDS"`
	KeycloakUrl       string   `mapstructure:"KEYCLOAK_URL"`
	OAuthClientSecret string   `mapstructure:"CLIENT_SECRET"`
	IataServiceUrl    string   `mapstructure:"IATA_SERVICE_URL"`
	OAuthClientId     string   `mapstructure:"OAUTH_CLIENT_ID"`
	OAuthUser         string   `mapstructure:"OAUTH_USER"`
	OAuthPassword     string   `mapstructure:"OAUTH_PASSWORD"`
}

func NewConfig() (error, *Config) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	env := Config{}
	err = viper.Unmarshal(&env)

	return err, &env
}
