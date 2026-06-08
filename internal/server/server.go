package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/files"
	"github.com/Azat201003/summorist-mores/lib/validators"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/Azat201003/summorist-shared/gen/go/users"
	"google.golang.org/grpc"
)

type moreServer struct {
	pb.UnimplementedMoresServer
	UsersClient users.UsersClient
	DMC         database.DatabaseMetasClient
}

func (s *moreServer) GetFiltered(request *pb.Meta, stream grpc.ServerStreamingServer[pb.Meta]) error {
	metas, err := s.DMC.RecieveFiltered(request)
	if err != nil {
		return err
	}

	log.Println(metas)

	for i := range metas {
		if err = stream.Send(&metas[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *moreServer) DownloadMore(request *pb.DownloadRequest, stream grpc.ServerStreamingServer[pb.Part]) error {
	metas, err := s.DMC.RecieveFiltered(&pb.Meta{MoreId: request.Data.MoreId})
	if err != nil {
		return err
	}
	if metas == nil {
		return ErrNotFound
	}

	var offset uint32 = 0
	buffer := make([]byte, request.Data.BlockSize)

	var i uint32 = 0
	for {
		n, err := files.ReadFile(uint32(request.Data.MoreId), offset, buffer)
		if err != nil {
			return err
		}
		var a []byte = make([]byte, n)
		copy(a, buffer)
		stream.Send(&pb.Part{Size: n, Data: a, Number: i})

		if n < request.Data.BlockSize {
			return nil
		}
		i++
		offset += n
	}
}

func (s *moreServer) UploadMore(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.Meta]) error {
	request, err := stream.Recv()
	var header *pb.ExchangeData
	if header = request.GetData(); header == nil {
		return ErrNoHeader
	}
	response, err := s.UsersClient.Authorize(context.Background(), &users.AuthRequest{JwtToken: header.JwtToken})
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return ErrNotAuthorized
	}
	metas, err := s.DMC.RecieveFiltered(&pb.Meta{MoreId: header.MoreId})
	if err != nil {
		return err
	}
	if metas == nil {
		return ErrNotFound
	}
	return nil
}

func NewServer() (*moreServer, error) {
	server := new(moreServer)
	return server, nil
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
