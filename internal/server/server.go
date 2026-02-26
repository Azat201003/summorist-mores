package server

import (
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	users "github.com/Azat201003/summorist-shared/gen/go/users"
	"google.golang.org/grpc"
	"net"
	"log"
	"os"
	"fmt"
)

type moreServer struct {
	pb.UnimplementedMoresServer
	usersClient *users.UsersClient
}

func (s *moreServer) UploadMore(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.Meta]) error {
		return nil
}

func newServer() *moreServer {
	server := new(moreServer)
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

