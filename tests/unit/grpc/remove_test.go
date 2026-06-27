package grpc

import (
	"context"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/files"
	"github.com/Azat201003/summorist-mores/internal/server"
	"github.com/Azat201003/summorist-mores/lib/client"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRemoveOk(t *testing.T) {
	// Preparing
	os.Setenv("MORES_HOST", "127.0.0.1")
	os.Setenv("MORES_PORT", "8005")
	os.Setenv("MORES_FILE_PREFIX", "../testdata/grpc_")
	os.Setenv("MORES_FILE_SUFFIX", ".txt")

	dbc := new(database.DatabaseClient)
	sqlDB, mock, err := dbc.GetMock()
	assert.NoError(t, err)
	defer sqlDB.Close()

	ms, err := server.NewServer()
	assert.NoError(t, err)
	ms.DMC = dbc
	ctx, cancel := context.WithCancel(t.Context())
	go ms.Start(ctx)
	defer cancel()

	conn, err := client.Connect()
	assert.NoError(t, err, "Cannot estabilish connection")
	defer conn.Close()
	client := pb.NewMoresClient(conn)

	moreId := uint64(2)
	userId := uint64(1)
	title := "title"

	files.RemoveFile(moreId)

	// Do
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "metas"`)).WillReturnRows(sqlmock.NewRows([]string{"more_id", "creator_id", "title"}).AddRow(moreId, userId, title))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "metas" WHERE id = $1`)).WithArgs(moreId).WillReturnResult(sqlmock.NewResult(int64(moreId), 1))
	mock.ExpectCommit()
	response, err := client.RemoveMore(t.Context(), &pb.RemoveRequest{MoreId: moreId})
	assert.NoError(t, err)

	// Check
	content, err := os.ReadFile(files.GetFilePath(moreId))
	log.Println(string(content))
	assert.Equal(t, title, response.Title)
	assert.Equal(t, userId, response.CreatorId)
	assert.Equal(t, moreId, response.MoreId)

	assert.NoError(t, mock.ExpectationsWereMet())
}
