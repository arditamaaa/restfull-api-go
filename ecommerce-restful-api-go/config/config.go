package config

import (
	"simple-api-go/util"

	"github.com/spf13/viper"
)

var (
	AppHost         string
	AppPort         int
	DBHost          string
	DBPort          int
	DBUsername      string
	DBPassword      string
	DBName          string
	TokenExpMinutes int
	SessionKey      string
)

func init() {
	loadConfig()

	// server configuration
	AppHost = viper.GetString("APP_HOST")
	AppPort = viper.GetInt("APP_PORT")

	// database configuration
	DBHost = viper.GetString("DB_HOST")
	DBPort = viper.GetInt("DB_PORT")
	DBUsername = viper.GetString("DB_USERNAME")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")

	// token configuration
	TokenExpMinutes = viper.GetInt("TOKEN_EXP_MINUTES")

	SessionKey = viper.GetString("SESSION_KEY")
}

func loadConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err == nil {
		util.Log.Infof("Config file loaded from %s", viper.ConfigFileUsed())
		return
	}
	util.Log.Error("Failed to load any config file")
}
