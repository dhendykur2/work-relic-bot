package app

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type IRepository interface {
}

type levelAppRepository struct {
	Conn *leveldb.DB
}

func NewLevelDBRepository(Conn *leveldb.DB) IRepository {
	return &levelAppRepository{Conn}
}
