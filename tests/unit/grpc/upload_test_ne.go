package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/server"
	"github.com/Azat201003/summorist-mores/lib/client"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/Azat201003/summorist-shared/gen/go/users"
	usersmock "github.com/Azat201003/summorist-shared/gen/go/users/mock"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	// Preparing
	os.Setenv("MORES_HOST", "127.0.0.1")
	os.Setenv("MORES_PORT", "8002")
	os.Setenv("MORES_FILE_PREFIX", "../testdata/grpc_")
	os.Setenv("MORES_FILE_SUFFIX", ".txt")

	dbc := new(database.DatabaseClient)
	sqlDB, mock, err := dbc.GetMock()
	assert.NoError(t, err)
	defer sqlDB.Close()

	usersMock := usersmock.NewMockUsersClient(gomock.NewController(t))

	ms, err := server.NewServer()
	assert.NoError(t, err)
	ms.DMC = dbc
	ms.UsersClient = usersMock
	ctx, cancel := context.WithCancel(t.Context())
	go ms.Start(ctx)
	defer cancel()

	conn, err := client.Connect()
	assert.NoError(t, err, "Cannot estabilish connection")
	defer conn.Close()
	client := pb.NewMoresClient(conn)

	jwt := "some jwt"
	moreId := uint64(1)
	userId := uint64(1)

	// Do
	usersMock.EXPECT().Authorize(t.Context(), &users.AuthRequest{JwtToken: jwt}).Return(userId, int32(0))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta"`)).WillReturnRows(sqlmock.NewRows([]string{"more_id", "creator_id", "title"}).AddRow(moreId, userId, "title"))
	stream, err := client.DownloadMore(t.Context(), &pb.DownloadRequest{Data: &pb.ExchangeData{MoreId: moreId, JwtToken: jwt, BlockSize: 8}})
	assert.NoError(t, err)

	result := []byte{}

	for {
		part, err := stream.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		result = append(result, part.Data...)
	}

	log.Println(string(result))

	// Check
	content, err := os.ReadFile(os.Getenv("MORES_FILE_PREFIX") + fmt.Sprint(moreId) + os.Getenv("MORES_FILE_SUFFIX"))
	assert.Equal(t, content, result)
}
