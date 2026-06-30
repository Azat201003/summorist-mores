package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseClient struct {
	DB *gorm.DB
}

func (dbc *DatabaseClient) Init() error {
	db, err := gorm.Open(postgres.Open(os.Getenv("MORES_POSTGRES_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}
	dbc.DB = db
	return nil
}

func (dbc *DatabaseClient) GetMock() (sqlDB *sql.DB, mock sqlmock.Sqlmock, err error) { // Instead of Init()
	sqlDB, mock, err = sqlmock.New()
	if err != nil {
		return
	}
	dbc.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL log
				Colorful:                  false,       // Disable color
			},
		),
	})
	return
}

func (dbc *DatabaseClient) RecieveFiltered(filter *pb.Filter) ([]*pb.Meta, error) {
	metas := []*pb.Meta{}
	var query = dbc.DB.Table("metas")
	if filter.MoreId != 0 {
		query = query.
			Where("more_id = ?", filter.MoreId)
	}
	if filter.CreatorId != 0 {
		query = query.
			Where("creator_id = ?", filter.CreatorId)
	}
	if filter.Query != "" {
		query = query.
    	Select("*, ts_rank(search_vector, plainto_tsquery('simple', ?)) AS rank", filter.Query).
    	Where("search_vector @@ plainto_tsquery('simple', ?)", filter.Query).
    	Order("rank DESC")
	}
	query = query.Find(&metas)
	return metas, query.Error
}

func (dbc *DatabaseClient) DeleteMore(id uint64) error {
	result := dbc.DB.Table("metas").Where("id = ?", id).Delete(&pb.Meta{})
	return result.Error
}

func (dbc *DatabaseClient) CreateMore(more *pb.Meta) (uint64, error) {
	result := dbc.DB.Table("metas").Create(more)
	return more.MoreId, result.Error
}

func (dbc *DatabaseClient) UpdateMore(more *pb.Meta) error {
	result := dbc.DB.Table("metas").Where("more_id = ?", more.MoreId).Updates(more)
	return result.Error
}
