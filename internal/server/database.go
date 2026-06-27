package server

import (
	"github.com/Azat201003/summorist-mores/internal/database"
	pb "github.com/Azat201003/summorist-shared/gen/go/mores"
)

func getMetaById(dmc database.DatabaseMetasClient, moreId uint64) (*pb.Meta, error) {
	metas, err := dmc.RecieveFiltered(&pb.Meta{MoreId: moreId})
	if err != nil {
		return nil, err
	}
	if metas == nil {
		return nil, ErrNotFound
	}
	return metas[0], nil
}
