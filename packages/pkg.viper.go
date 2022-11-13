package packages

import (
	"os"

	"github.com/spf13/viper"
)

func ViperLoadConfig() error {
	if _, ok := os.LookupEnv("GO_ENV"); !ok {
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		return err
	}

	if env := os.Getenv("GO_ENV"); env == "development" {
		viper.Debug()
	}

	return nil
}

func GetString(name string) string {
	if _, ok := os.LookupEnv("GO_ENV"); !ok {
		return viper.GetString(name)

	}
	return viper.GetString(name)
}
