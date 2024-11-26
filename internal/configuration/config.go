package configuration

import (
	"github.com/spf13/viper"
)

type Config struct {
	GcProjectId    string   `mapstructure:"GCLOUD_PROJECT_ID"`
	GcLocation     string   `mapstructure:"GCLOUD_LOCATION"`
	GcBucketName   string   `mapstructure:"GCLOUD_BUCKETNAME"`
	AiModel        string   `mapstructure:"AI_MODEL"`
	OAuthIssuer    string   `mapstructure:"OAUTH_ISSUER"`
	OAuthClientIds []string `mapstructure:"OAUTH_CLIENT_IDS"`
}

func NewConfig() (error, *Config) {
	env := Config{}
	v := viper.New()
	v.SetConfigFile(".env")
	//docker
	v.AddConfigPath("/app")
	//local
	v.AddConfigPath("./")

	err := v.ReadInConfig()
	if err != nil {
		return err, nil
	}

	err = v.Unmarshal(&env)
	if err != nil {
		return err, nil
	}
	return nil, &env
}
