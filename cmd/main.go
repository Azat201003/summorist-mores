package main

import (
	"github.com/Azat201003/summorist-mores/internal/config"
)

func main() {
	conf := config.GetConfig()
	if !conf.ConfigIncluded {
		panic("Config's not included.")
	}
	if !conf.SecretsIncluded {
		panic("Secrets' not included.")
	}
}

