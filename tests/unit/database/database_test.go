package database_test

import (
	"reflect"
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
	field := reflect.ValueOf(dbc).Elem().FieldByName("DB")
	field.Set(reflect.ValueOf(gdb))

	return dbc, mock
}

func TestRecieveFiltered(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	filter := &pb.Meta{MoreId: 1}
	rows := sqlmock.NewRows([]string{"more_id"}).AddRow(1)
	mock.ExpectQuery(`SELECT \* FROM "meta" WHERE "meta"\."more_id" = \$1`).WillReturnRows(rows)

	metas, err := dbc.RecieveFiltered(filter)
	assert.NoError(t, err)
	assert.Len(t, *metas, 1)
	assert.Equal(t, uint64(1), (*metas)[0].MoreId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMore(t *testing.T) {
	dbc, mock := setupMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "meta" WHERE id = \$1`).WillReturnResult(sqlmock.NewResult(0, 1))
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
	mock.ExpectExec(`INSERT INTO "meta"`).WillReturnResult(sqlmock.NewResult(1, 1))
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
	mock.ExpectExec(`UPDATE "meta" SET "more_id"=\$1 WHERE more_id = \$2`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dbc.UpdateMore(more)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
