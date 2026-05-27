package grpc

import (
	"context"
	"io"
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
	dbc := new(database.DatabaseClient)
	mock, err := dbc.GetMock()
	assert.NoError(t, err)

	ms := server.NewServer()
	ms.DMC = dbc
	ctx, cancel := context.WithCancel(t.Context())
	go ms.Start(ctx)
	defer cancel()

	conn, err := client.Connect()
	assert.NoError(t, err, "Cannot estabilish connection")
	defer conn.Close()

	client := pb.NewMoresClient(conn)

	// Do
	mock.ExpectQuery("*").WillReturnRows(sqlmock.NewRows([]string{"1", "title", "2"}))
	ans, err := client.GetFiltered(t.Context(), &pb.Meta{})

	// Check
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	meta, err := ans.Recv()
	assert.NoError(t, err)
	assert.Equal(t, meta.GetCreatorId(), uint64(2))
	assert.Equal(t, meta.GetMoreId(), uint64(1))
	assert.Equal(t, meta.GetTitle(), "title")

	meta, err = ans.Recv()
	assert.ErrorIs(t, err, io.EOF)
}
