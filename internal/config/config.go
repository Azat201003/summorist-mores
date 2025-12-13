package config

import (
	"os"
)

type Config struct {
	FilePrefix		 string
	FilePostfix 	 string
	DBHost			 string
	DBPort			 string
	DBUser			 string
	DBName			 string
	DBPassword		 string
	ConfigIncluded	 bool
	SecretsIncluded	 bool
}

var config Config
var isRead = false

func GetConfig() Config {
	if !isRead {
		config.ConfigIncluded		= os.Getenv("CONFIG_INCLUDED") != ""
		config.SecretsIncluded		= os.Getenv("SECRETS_INCLUDED") != ""

		config.FilePrefix  		   	= os.Getenv("FILE_PREFIX")
		config.FilePostfix			= os.Getenv("FILE_POSTFIX")

		config.DBHost 	   			= os.Getenv("POSTGRES_HOST")
		config.DBPort 	   			= os.Getenv("POSTGRES_PORT")
		config.DBUser	   			= os.Getenv("POSTGRES_USER")
		config.DBName 	   			= os.Getenv("POSTGRES_DB")
		config.DBPassword  			= os.Getenv("POSTGRES_PASSWORD")
	
		isRead = true
	}
	return config
}
