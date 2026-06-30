package database_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Azat201003/summorist-mores/internal/database"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
)

func setupMockDB(t *testing.T) (*database.DatabaseClient, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		DriverName:           "postgres",
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	dbc := &database.DatabaseClient{}
	dbc.DB = gdb

	return dbc, mock
}

func TestRecieveFiltered_ByMoreId(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"more_id", "descripiton"}).AddRow(1, "")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "metas" WHERE more_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	metas, err := dbc.RecieveFiltered(&pb.Filter{MoreId: 1})
	assert.NoError(t, err)
	assert.Len(t, metas, 1)
	assert.Equal(t, uint64(1), metas[0].MoreId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecieveFiltered_ByQuery(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"more_id", "descripiton", "rank"}).AddRow(1, "", 0.5)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT *, ts_rank(search_vector, plainto_tsquery('simple', $1)) AS rank FROM "metas" WHERE search_vector @@ plainto_tsquery('simple', $2) ORDER BY rank DESC`)).
		WithArgs("test", "test").
		WillReturnRows(rows)

	metas, err := dbc.RecieveFiltered(&pb.Filter{Query: "test"})
	assert.NoError(t, err)
	assert.Len(t, metas, 1)
	assert.Equal(t, uint64(1), metas[0].MoreId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecieveFiltered_ByMoreIdAndQuery(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"more_id", "descripiton", "rank"}).AddRow(1, "", 0.5)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT *, ts_rank(search_vector, plainto_tsquery('simple', $1)) AS rank FROM "metas" WHERE more_id = $2 AND search_vector @@ plainto_tsquery('simple', $3) ORDER BY rank DESC`)).
		WithArgs("test", 1, "test").
		WillReturnRows(rows)

	metas, err := dbc.RecieveFiltered(&pb.Filter{MoreId: 1, Query: "test"})
	assert.NoError(t, err)
	assert.Len(t, metas, 1)
	assert.Equal(t, uint64(1), metas[0].MoreId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecieveFiltered_NoFilters(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"more_id", "descripiton"}).AddRow(1, "")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "metas"`)).
		WillReturnRows(rows)

	metas, err := dbc.RecieveFiltered(&pb.Filter{})
	assert.NoError(t, err)
	assert.Len(t, metas, 1)
	assert.Equal(t, uint64(1), metas[0].MoreId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMore(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "metas" WHERE id = \$1`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dbc.DeleteMore(1)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateMore(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	more := &pb.Meta{MoreId: 1}
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "metas"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := dbc.CreateMore(more)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateMore(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	more := &pb.Meta{MoreId: 1}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "metas" SET "more_id"=\$1 WHERE more_id = \$2`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dbc.UpdateMore(more)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
