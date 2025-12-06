package config

import (
	"os"
)

type Config struct {
	FilePrefix	 string
	FilePostfix  string
}

var config Config
var isRead = false

func GetConfig() Config {
	if !isRead {
		// config.env.*
		config.FilePrefix  = os.Getenv("FILE_PREFIX")
		config.FilePostfix = os.Getenv("FILE_POSTFIX")
	
		isRead = true
	}
	return config
}
