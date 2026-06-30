package server

import (
	"context"
	"fmt"
	"io"
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

func (ms *moreServer) GetFiltered(request *pb.Filter, stream grpc.ServerStreamingServer[pb.Meta]) error {
	metas, err := ms.DMC.RecieveFiltered(request)
	if err != nil {
		return err
	}

	log.Println(metas)

	for i := range metas {
		if err = stream.Send(metas[i]); err != nil {
			return err
		}
	}
	return nil
}

func (ms *moreServer) DownloadMore(request *pb.DownloadRequest, stream grpc.ServerStreamingServer[pb.Part]) error {
	_, err := getMetaById(ms.DMC, request.Data.MoreId)
	if err != nil {
		return err
	}

	var offset uint32 = 0
	buffer := make([]byte, request.Data.BlockSize)

	var i uint32 = 0
	for {
		n, err := files.ReadFile(request.Data.MoreId, offset, buffer)
		if err != nil {
			return err
		}
		a := make([]byte, n)
		copy(a, buffer)
		stream.Send(&pb.Part{Size: n, Data: a, Number: i})

		if n < request.Data.BlockSize {
			return nil
		}
		i++
		offset += n
	}
}

func (ms *moreServer) CreateMore(ctx context.Context, request *pb.CreateRequest) (*pb.Meta, error) {
	response, err := ms.UsersClient.Authorize(context.Background(), &users.AuthRequest{JwtToken: request.JwtToken})
	if err != nil {
		return nil, ErrNotAuthorized
	}
	id, err := ms.DMC.CreateMore(&pb.Meta{CreatorId: response.UserId, Title: request.Title, Descripiton: request.Descripiton})
	log.Println(id)
	if err != nil {
		return nil, err
	}
	more, err := getMetaById(ms.DMC, id)
	if err != nil {
		return nil, err
	}
	log.Println(more)
	return more, err
}

func (ms *moreServer) UploadMore(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.Meta]) error {
	request, err := stream.Recv()
	if err != nil {
		return err
	}
	var header *pb.ExchangeData
	if header = request.GetData(); header == nil {
		return ErrNoHeader
	}
	authedUserId, err := authorize(ms.UsersClient, header.JwtToken)
	if err != nil {
		return err
	}
	more, err := getMetaById(ms.DMC, header.MoreId)
	if err != nil {
		return err
	}
	if more.CreatorId != authedUserId {
		return ErrNotPermitted
	}

	for {
		request, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		part := request.GetPart()
		if part == nil {
			return ErrNoContent
		}

		files.WriteFile(more.MoreId, part.Data)
	}

	return stream.SendAndClose(more)
}

func (ms *moreServer) RemoveMore(ctx context.Context, request *pb.RemoveRequest) (response *pb.Meta, err error) {
	metas, err := ms.DMC.RecieveFiltered(&pb.Filter{MoreId: request.MoreId})
	if err != nil {
		return
	}
	if len(metas) == 0 {
		err = ErrNotFound
		return
	}
	response = metas[0]
	err = ms.DMC.DeleteMore(request.MoreId)
	return
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
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			log.Println("Server stopped")
			s.Stop()
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func NewServer() (*moreServer, error) {
	server := new(moreServer)
	return server, nil
}
