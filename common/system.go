package common

import (
	"database/sql"

	"github.com/hiddedorhout/travel/services/travel"

	"github.com/hiddedorhout/travel/services/kv_store"

	_ "github.com/mattn/go-sqlite3"
)

// System is the base system service type
type System struct {
	db      *sql.DB
	kvStore *kv_store.Service
	travels *travel.Service
}

// New initiates the system
func New(dbName string) (*System, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	kvStore, err := kv_store.New(db)
	if err != nil {
		return nil, err
	}

	travel, err := travel.New(db)
	if err != nil {
		return nil, err
	}

	return &System{
		db:      db,
		kvStore: kvStore,
		travels: travel,
	}, nil
}
