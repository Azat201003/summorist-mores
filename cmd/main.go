package main

import (
	"context"

	"github.com/Azat201003/summorist-mores/internal/server"
)

func main() {
	server.StartServer(context.Background())
}
