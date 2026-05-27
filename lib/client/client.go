package client

import (
	"fmt"
	"os"

	"github.com/Azat201003/summorist-mores/lib/validators"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect() (*grpc.ClientConn, error) {
	host := os.Getenv("MORES_HOST")
	if err := validators.ValidateHost(host); err != nil {
		return nil, fmt.Errorf(`Host (environment variable "MORES_HOST") is invalid: %v`, err)
	}

	port := os.Getenv("MORES_PORT")

	if err := validators.ValidatePort(port); err != nil {
		return nil, fmt.Errorf(`Port (environment variable "MORES_PORT") is invalid: %v`, err)
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
