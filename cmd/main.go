package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/server"
	"github.com/Azat201003/summorist-shared/gen/go/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	host := os.Getenv("USERS_HOST")
	port := os.Getenv("USERS_PORT")

	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	s.UsersClient = users.NewUsersClient(conn)

	s.Start(context.Background())
}
