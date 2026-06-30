package grpc

import (
	"context"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/server"
	"github.com/Azat201003/summorist-mores/lib/client"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetFilteredOk(t *testing.T) {
	// Preparing
	os.Setenv("MORES_HOST", "127.0.0.1")
	os.Setenv("MORES_PORT", "8001")

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

	const moreId uint64 = 1
	const creatorId uint64 = 1
	const title = "title"
	const description = ""

	// Do
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "metas"`)).
		WillReturnRows(sqlmock.NewRows([]string{"more_id", "creator_id", "title", "descripiton"}).AddRow(moreId, creatorId, title, description))
	ans, err := client.GetFiltered(t.Context(), &pb.Filter{})

	// Check
	assert.NoError(t, err)

	meta, err := ans.Recv()
	assert.NoError(t, err)

	assert.Equal(t, creatorId, meta.GetCreatorId())
	assert.Equal(t, moreId, meta.GetMoreId())
	assert.Equal(t, title, meta.GetTitle())
	assert.Equal(t, description, meta.GetDescripiton())

	assert.NoError(t, mock.ExpectationsWereMet())

	meta, err = ans.Recv()
	assert.ErrorIs(t, io.EOF, err)
}
