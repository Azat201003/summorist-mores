package database

import (
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
)

type DatabaseMetasClient interface {
	RecieveFiltered(filter *pb.Meta) (*[]pb.Meta, error)
	DeleteMore(id uint64) error
	CreateMore(more *pb.Meta) (uint64, error)
	UpdateMore(more *pb.Meta) error
}
