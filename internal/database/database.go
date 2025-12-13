package database

import (
	"fmt"
	"github.com/Azat201003/summorist-mores/internal/config"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseClient struct {
	DB *gorm.DB
}

func (dbc *DatabaseClient) Init() {
	conf := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai",
		conf.DBHost,
		conf.DBUser,
		conf.DBPassword,
		conf.DBName,
		conf.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err) // or handle error
	}
	dbc.DB = db
}

func (dbc *DatabaseClient) RecieveFiltered(filter *pb.Meta) (*[]pb.Meta, error) {
	metas := new([]pb.Meta)
	result := dbc.DB.Find(metas, filter)
	return metas, result.Error
}

func (dbc *DatabaseClient) DeleteMore(id uint64) error {
	result := dbc.DB.Where("id = ?", id).Delete(&pb.Meta{})
	return result.Error
}

func (dbc *DatabaseClient) CreateMore(more *pb.Meta) (uint64, error) {
	result := dbc.DB.Create(more)
	return more.MoreId, result.Error
}

func (dbc *DatabaseClient) UpdateMore(more *pb.Meta) error {
	result := dbc.DB.Where("more_id = ?", more.MoreId).Updates(more)
	return result.Error
}
