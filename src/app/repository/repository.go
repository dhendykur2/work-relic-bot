package app

import (
	"encoding/json"
	"work-relic-bot/src/models"

	"github.com/syndtr/goleveldb/leveldb"
)

type IRepository interface {
	Insert(val json.RawMessage, username string) error
	GetUndone(username string) ([]models.Task, error)
}

type levelAppRepository struct {
	Conn *leveldb.DB
}

func NewLevelDBRepository(Conn *leveldb.DB) IRepository {
	return &levelAppRepository{Conn}
}

func (l *levelAppRepository) GetUndone(username string) ([]models.Task, error) {
	data, err := l.Conn.Get([]byte(username+"_undone"), nil)
	if err != nil {
		return nil, err
	}
	var result []models.Task
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (l *levelAppRepository) Insert(val json.RawMessage, username string) error {
	err := l.Conn.Put([]byte(username+"_undone"), val, nil)
	if err != nil {
		return err
	}
	return nil
}
