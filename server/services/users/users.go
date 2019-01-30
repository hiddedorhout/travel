package users

import (
	"database/sql"
	"encoding/base64"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db              *sql.DB
	createUserStmnt *sql.Stmt
	getUserStmnt    *sql.Stmt
}

func New(db *sql.DB) (*Service, error) {

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(
		ID NOT NULL TET PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`); err != nil {
		return nil, err
	}

	createUserStmnt, err := db.Prepare(`INSERT INTO users (ID, username, password) VALUES (?,?,?)`)
	if err != nil {
		return nil, err
	}

	getUserStmnt, err := db.Prepare(`SELECT ID, password FROM users WHERE username = ?`)
	if err != nil {
		return nil, err
	}

	return &Service{
		db:              db,
		createUserStmnt: createUserStmnt,
		getUserStmnt:    getUserStmnt,
	}, nil
}

// RegisterUser creates a digest over
func (s *Service) RegisterUser(username, plainTextPassword string) (*string, error) {

	id := uuid.New().String()

	encryotedPwd, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 10)
	if err != nil {
		return nil, err
	}
	b64Pwd := base64.StdEncoding.EncodeToString(encryotedPwd)

	if _, err := s.createUserStmnt.Exec(id, username, b64Pwd); err != nil {
		return nil, err
	}

	return &id, nil
}

// CheckPassword authenticates the user and return true if user/password is ok
func (s *Service) CheckPassword(username, password string) (*string, error) {
	var b64EncodedPwd string
	var ID string
	if err := s.getUserStmnt.QueryRow(username).Scan(&ID, &b64EncodedPwd); err != nil {
		return nil, err
	}

	rawPwdHash, err := base64.StdEncoding.DecodeString(b64EncodedPwd)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), rawPwdHash); err != nil {
		return nil, err
	}

	return &ID, nil
}
