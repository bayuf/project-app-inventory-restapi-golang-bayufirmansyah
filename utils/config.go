package utils

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	JWTKey      string
	AppName     string
	Port        string
	Debug       bool
	PageLimit   int
	PathLogging string
	DB          DatabaseConfiguration
}

type DatabaseConfiguration struct {
	DBName   string
	UserName string
	Password string
	HostName string
	SSLMode  string
}

func ReadConfiguration() (*Configuration, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Configuration{
		JWTKey:      viper.GetString("JWT_KEY"),
		AppName:     viper.GetString("APP_NAME"),
		Port:        viper.GetString("PORT"),
		Debug:       viper.GetBool("DEBUG"),
		PageLimit:   viper.GetInt("PAGE_LIMIT"),
		PathLogging: viper.GetString("PATH_LOGGING"),
		DB: DatabaseConfiguration{
			DBName:   viper.GetString("DB_NAME"),
			UserName: viper.GetString("DB_USERNAME"),
			Password: viper.GetString("DB_PASSWORD"),
			HostName: viper.GetString("DB_HOSTNAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
	}, nil
}
