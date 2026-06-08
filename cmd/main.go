package main

import (
	"context"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/server"
)

func main() {
	s, err := server.NewServer()
	dmc := database.DatabaseClient{}
	if err != nil {
		panic(err)
	}
	err = dmc.Init()
	s.DMC = &dmc
	if err != nil {
		panic(err)
	}
	s.Start(context.Background())
}
