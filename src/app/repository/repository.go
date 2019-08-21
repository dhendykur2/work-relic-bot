package app

import (
	"encoding/json"
	"work-relic-bot/src/models"

	"github.com/syndtr/goleveldb/leveldb"
)

// IRepository to register all repository function
type IRepository interface {
	Put(key string, val json.RawMessage) error
	Get(key string) ([]models.Task, error)
}

type levelAppRepository struct {
	Conn *leveldb.DB
}

// NewLevelDBRepository to register DB into repository
func NewLevelDBRepository(Conn *leveldb.DB) IRepository {
	return &levelAppRepository{Conn}
}

func (l *levelAppRepository) Get(key string) ([]models.Task, error) {
	data, err := l.Conn.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	var result []models.Task
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (l *levelAppRepository) Put(key string, val json.RawMessage) error {
	err := l.Conn.Put([]byte(key), val, nil)
	if err != nil {
		return err
	}
	return nil
}
