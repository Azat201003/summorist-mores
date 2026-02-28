package server_tests

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	users_mock "github.com/Azat201003/summorist-shared/gen/go/users/mock"

	"github.com/Azat201003/summorist-mores/internal/database"
	"github.com/Azat201003/summorist-mores/internal/server"
)

type serverSuite struct {
	suite.Suite
	moresClient *pb.MoresClient
	usersClientMock *users_mock.MockUsersClient
	dbc         *database.DatabaseClient
	dbmock      sqlmock.Sqlmock
	lis 				net.Listener
	db 					*sql.DB
	gomockCtrl *gomock.Controller
}

func (s *serverSuite) SetupSuite() {
	// Mock database
	log.Println("1. Mock db setting up")
	db, mock, err := sqlmock.New()
	s.db = db
	s.NoError(err)
	s.dbmock = mock

	// Database
	log.Println("2. Database setting up")
	dialector := postgres.New(postgres.Config{
		Conn: db,
		DriverName: "postgres",
	})
	gormdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	s.NoError(err)
	s.dbc = &database.DatabaseClient{DB: gormdb}
	
	// User service mock
	ctrl := gomock.NewController(s.T())
	s.gomockCtrl = ctrl
	s.usersClientMock = users_mock.NewMockUsersClient(ctrl)

	// This service server
	log.Println("3. Service server setting up")
	s.lis, _ = net.Listen("tcp", fmt.Sprintf("%v:%v", "0.0.0.0", os.Getenv("MORES_PORT")))
	grpcServer := grpc.NewServer()
	pb.RegisterMoresServer(grpcServer, &server.MoreServer{DBC: s.dbc, UsersClient: s.usersClientMock})
	go grpcServer.Serve(s.lis)
	
	// This service client
	log.Println("4. Service client setting up")
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", "0.0.0.0", os.Getenv("USERS_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))

	s.NoError(err)

	client := pb.NewMoresClient(conn)
	s.moresClient = &client

	// End
	log.Println("5. All was set up")
}

func (s *serverSuite) TearDownSuite() {
	s.lis.Close()
	s.db.Close()
	s.gomockCtrl.Finish()
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverSuite))
}
