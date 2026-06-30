package grpc

import (
	"context"
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

func TestCreateMoreOk(t *testing.T) {
	// Prepare
	os.Setenv("MORES_HOST", "127.0.0.1")
	os.Setenv("MORES_PORT", "8003")
	os.Setenv("MORES_FILE_PREFIX", "../grpc_")
	os.Setenv("MORES_FILE_SUFFIX", ".txt")

	dbc := new(database.DatabaseClient)
	sqlDB, mock, err := dbc.GetMock()
	assert.NoError(t, err)
	defer sqlDB.Close()

	usersMock := usersmock.NewMockUsersClient(gomock.NewController(t))

	server, err := server.NewServer()
	server.DMC = dbc
	server.UsersClient = usersMock
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	go server.Start(ctx)
	defer cancel()

	conn, err := client.Connect()
	assert.NoError(t, err)
	defer conn.Close()
	client := pb.NewMoresClient(conn)

	const jwt = "some jwt"
	const moreId uint64 = 2
	const userId uint64 = 3
	const title = "Something wonderful"
	const description = "Read title"

	// Do
	usersMock.EXPECT().Authorize(context.Background(), &users.AuthRequest{JwtToken: jwt}).Return(&users.AuthResponse{UserId: userId}, nil)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "metas" ("title","descripiton","creator_id","more_id") VALUES ($1,$2,$3,$4)`)).
		WithArgs(title, description, userId, 0).
		WillReturnResult(sqlmock.NewResult(int64(moreId), 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "metas"`)).
		WillReturnRows(sqlmock.NewRows([]string{"more_id", "creator_id", "title", "descripiton"}).AddRow(moreId, userId, title, description))
	response, err := client.CreateMore(context.Background(), &pb.CreateRequest{Title: title, Descripiton: description, JwtToken: jwt})

	// Check
	assert.NoError(t, err)
	assert.Equal(t, moreId, response.MoreId)
	assert.Equal(t, userId, response.CreatorId)
	assert.Equal(t, title, response.Title)
	assert.Equal(t, description, response.Descripiton)
	assert.NoError(t, mock.ExpectationsWereMet())
}
