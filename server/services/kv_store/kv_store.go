package kv_store

import "database/sql"

// Service is a key value store servive type
type Service struct {
	db         *sql.DB
	storeStmnt *sql.Stmt
	getStmnt   *sql.Stmt
}

// New initiates a new key_value store service
func New(db *sql.DB) (*Service, error) {

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS keyvalues(
		ID INTEGER NOT NULL PRIMARY KEY,
		key TEXT UNIQUE NOT NULL,
		value TEXT
	)`); err != nil {
		return nil, err
	}

	storeStmnt, err := db.Prepare(`INSERT INTO keyvalues (key, value) VALUES (?, ?)`)
	if err != nil {
		return nil, err
	}

	getStmnt, err := db.Prepare(`SELECT value FROM keyvalues WHERE key = ?`)
	if err != nil {
		return nil, err
	}

	return &Service{
		db:         db,
		getStmnt:   getStmnt,
		storeStmnt: storeStmnt,
	}, nil
}

// Store stores a key value string pair
func (s *Service) Store(key, value string) error {
	if _, err := s.storeStmnt.Exec(key, value); err != nil {
		return err
	}
	return nil
}

// Get returns the value (string) of a key
func (s *Service) Get(key string) (value *string, err error) {
	var result string
	if err := s.getStmnt.QueryRow(key).Scan(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
