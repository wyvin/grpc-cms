package setting

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func Setup() {
	var err error
	configType := "toml"
	configPath := "conf"

	env := os.Getenv("GRPC_ENV")
	if env != "production" {
		env = "default"
	}

	viper.SetConfigName(env)
	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("viper ReadInConfig err: %v", err)
	}
}
