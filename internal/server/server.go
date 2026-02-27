package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Azat201003/summorist-mores/internal/database"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	users "github.com/Azat201003/summorist-shared/gen/go/users"
	"google.golang.org/grpc"
)

type MoreServer struct {
	pb.UnimplementedMoresServer
	UsersClient *users.UsersClient
	DBC *database.DatabaseClient
}

func (s *MoreServer) UploadMore(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.Meta]) error {
		return nil
}

func newServer() *MoreServer {
	server := new(MoreServer)
	return server
}

func StartServer() {
	server := newServer()
	
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", os.Getenv("MORES_HOST"), os.Getenv("MORES_PORT")))
	if err != nil {
			log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMoresServer(s, server)

	log.Printf("Server listening on :%v\n", os.Getenv("MORES_PORT"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

