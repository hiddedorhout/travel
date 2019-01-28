package users

import (
	"database/sql"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db              *sql.DB
	createUserStmnt *sql.Stmt
	getUserStmnt    *sql.Stmt
}

func New(db *sql.DB) (*Service, error) {

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(
		ID NOT NULL INTEGER PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	)`); err != nil {
		return nil, err
	}

	createUserStmnt, err := db.Prepare(`INSERT INTO users (username, password) VALUES (?,?)`)
	if err != nil {
		return nil, err
	}

	getUserStmnt, err := db.Prepare(`SELECT password FROM users WHERE username = ?`)
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
func (s *Service) RegisterUser(username, plainTextPassword string) error {

	encryotedPwd, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 10)
	if err != nil {
		return err
	}
	b64Pwd := base64.StdEncoding.EncodeToString(encryotedPwd)

	if _, err := s.createUserStmnt.Exec(username, b64Pwd); err != nil {
		return err
	}

	return nil
}

// CheckPassword authenticates the user and return true if user/password is ok
func (s *Service) CheckPassword(username, password string) (bool, error) {
	var b64EncodedPwd string
	if err := s.getUserStmnt.QueryRow(username).Scan(&b64EncodedPwd); err != nil {
		return false, err
	}

	rawPwdHash, err := base64.StdEncoding.DecodeString(b64EncodedPwd)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), rawPwdHash); err != nil {
		return false, err
	}

	return true, nil
}
