package server_tests

import (
	"context"

	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/DATA-DOG/go-sqlmock"
)

func (s *serverSuite) TestRemoveOk() {
	moreId := uint64(1)

	s.dbmock.ExpectBegin()
	s.dbmock.ExpectExec("UPDATE").
		WithArgs(moreId, true).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.dbmock.ExpectCommit()
	
	(*s.moresClient).RemoveMore(context.Background(), &pb.RemoveRequest{
		JwtToken: "",
		MoreId: moreId,
	})

	s.NoError(s.dbmock.ExpectationsWereMet())
}

