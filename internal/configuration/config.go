package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"path/filepath"
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	env := Config{}
	v := viper.New()
	v.AutomaticEnv()
	if len(v.AllKeys()) == 0 {
		v.SetConfigType("env")
		v.SetConfigFile(".env")
		v.AddConfigPath(".")
		err := v.ReadInConfig()
		if err != nil {
			return err, nil
		}
	}
	slog.Info(v.ConfigFileUsed())
	err = v.Unmarshal(&env)
	if err != nil {
		return err, nil
	}
	return nil, &env
}
