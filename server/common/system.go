package common

import (
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"log"

	"github.com/hiddedorhout/travel/server/services/sessions"

	"github.com/hiddedorhout/travel/server/services/serviceCert"

	"github.com/hiddedorhout/travel/server/services/kv_store"
	"github.com/hiddedorhout/travel/server/services/travel"
	"github.com/hiddedorhout/travel/server/services/users"

	_ "github.com/mattn/go-sqlite3"
)

// System is the base system service type
type System struct {
	db       *sql.DB
	kvStore  *kv_store.Service
	users    *users.Service
	travels  *travel.Service
	sessions *sessions.Service
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

	users, err := users.New(db)
	if err != nil {
		return nil, err
	}

	log.Println("Generating service private key and cert")
	pkey, err := serviceCert.GenPkey()
	if err != nil {
		return nil, err
	}

	pkcs1 := x509.MarshalPKCS1PrivateKey(pkey)
	if err != nil {
		return nil, err
	}

	if err := kvStore.Store("pkey", base64.StdEncoding.EncodeToString(pkcs1)); err != nil {
		return nil, err
	}

	cert, err := serviceCert.SelfSignedCert(*pkey)
	if err != nil {
		return nil, err
	}
	if err := kvStore.Store("cert", base64.StdEncoding.EncodeToString(*cert)); err != nil {
		return nil, err
	}

	sessions, err := sessions.New(db, kvStore)
	if err != nil {
		return nil, err
	}

	return &System{
		db:       db,
		kvStore:  kvStore,
		travels:  travel,
		users:    users,
		sessions: sessions,
	}, nil
}
