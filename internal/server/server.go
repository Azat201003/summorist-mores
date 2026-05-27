package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/lib/validators"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/Azat201003/summorist-shared/gen/go/users"
	"google.golang.org/grpc"
)

type moreServer struct {
	pb.UnimplementedMoresServer
	usersClient *users.UsersClient
	DMC         database.DatabaseMetasClient
}

func (s *moreServer) UploadMore(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.Meta]) error {
	return nil
}

func NewServer() *moreServer {
	server := new(moreServer)
	return server
}

func (ms *moreServer) Start(ctx context.Context) {
	host := os.Getenv("MORES_HOST")
	if err := validators.ValidateHost(host); err != nil {
		log.Fatalf(`Host (environment variable "MORES_HOST") is invalid: %v`, err)
	}

	port := os.Getenv("MORES_PORT")

	if err := validators.ValidatePort(port); err != nil {
		log.Fatalf(`Port (environment variable "MORES_PORT") is invalid: %v`, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMoresServer(s, ms)

	log.Printf("Server listening on %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			s.Stop()
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
