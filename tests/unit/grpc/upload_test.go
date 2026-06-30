package grpc

import (
	"context"
	"io"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/files"
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
	// Preparing
	os.Setenv("MORES_HOST", "127.0.0.1")
	os.Setenv("MORES_PORT", "8004")
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

	const jwt = "some jwt"
	const moreId uint64 = 2
	const userId uint64 = 1
	const blockSize uint32 = 5
	const title = "title"
	const description = ""

	files.RemoveFile(moreId)

	// Do
	usersMock.EXPECT().Authorize(context.Background(), &users.AuthRequest{JwtToken: jwt}).Return(&users.AuthResponse{UserId: userId, Code: int32(0)}, nil)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "metas" WHERE more_id = $1`)).
		WithArgs(moreId).
		WillReturnRows(sqlmock.NewRows([]string{"more_id", "creator_id", "title", "descripiton"}).AddRow(moreId, userId, title, description))
	stream, err := client.UploadMore(t.Context())
	assert.NoError(t, err)
	err = stream.Send(&pb.UploadRequest{Request: &pb.UploadRequest_Data{Data: &pb.ExchangeData{JwtToken: jwt, MoreId: moreId, BlockSize: blockSize}}})
	assert.NoError(t, err)

	data := []byte("This is super mega ultra text, that is super ultra mega, also it's text!")
	offset := uint32(0)

	var more *pb.Meta

	for {
		err = stream.Send(&pb.UploadRequest{Request: &pb.UploadRequest_Part{Part: &pb.Part{Data: data[offset:min(offset+blockSize, uint32(len(data)))]}}})
		offset += blockSize
		assert.NoError(t, err)
		if offset >= uint32(len(data)) {
			stream.CloseSend()
			more, err = stream.CloseAndRecv()
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
			break
		}
	}

	// Check
	content, err := os.ReadFile(files.GetFilePath(moreId))
	log.Println(string(content))
	assert.Equal(t, data, content)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, moreId, more.MoreId)
	assert.Equal(t, title, more.Title)
	assert.Equal(t, description, more.Descripiton)
}
