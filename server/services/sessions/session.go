package sessions

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/hiddedorhout/travel/server/services/kv_store"
)

type Service struct {
	db                *sql.DB
	storeSessionStmnt *sql.Stmt
	kvStore           *kv_store.Service
}

func New(db *sql.DB, kvStore *kv_store.Service) (*Service, error) {

	return &Service{
		db:      db,
		kvStore: kvStore,
	}, nil
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtPayload struct {
	Iss string `json:"iss"`
	Aud string `json:"aud"`
	Sub string `json:"sub"`
	Iat int    `json:"iat`
	Exp int    `json:"exp"`
}

func (s *Service) GenerateSession(aud string) (*string, error) {
	header := jwtHeader{
		Alg: "RSA1_5",
		Typ: "JWT",
	}

	payload := jwtPayload{
		Iss: "Travel application",
		Aud: aud,
		Sub: "SessionToken",
		Iat: time.Now().Nanosecond(),
		Exp: time.Now().Add(60 * time.Second).Nanosecond(),
	}

	bpayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	bheader, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	encPayload := encode(bpayload)
	encHeader := encode(bheader)

	tbs := encHeader + "." + encPayload

	b64pkey, err := s.kvStore.Get("pkey")
	if err != nil {
		return nil, err
	}

	rawPkey, err := base64.StdEncoding.DecodeString(*b64pkey)
	if err != nil {
		return nil, err
	}
	pkey, err := x509.ParsePKCS1PrivateKey(rawPkey)
	if err != nil {
		return nil, err
	}

	signature, err := sign(pkey, []byte(tbs))
	if err != nil {
		return nil, err
	}

	b64sig := base64.URLEncoding.EncodeToString(*signature)

	jwt := tbs + "." + b64sig
	return &jwt, nil

}
func generateSession(aud string, pkey *rsa.PrivateKey) (*string, error) {
	header := jwtHeader{
		Alg: "RS256",
		Typ: "JWT",
	}

	payload := jwtPayload{
		Iss: "Travel application",
		Aud: aud,
		Sub: "SessionToken",
		Iat: time.Now().Nanosecond(),
		Exp: time.Now().Add(60 * time.Second).Nanosecond(),
	}

	bpayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	bheader, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	encPayload := encode(bpayload)
	encHeader := encode(bheader)

	tbs := encHeader + "." + encPayload

	h := sha256.New()
	h.Write([]byte(tbs))
	hashed := h.Sum(nil)

	signature, err := sign(pkey, hashed)
	if err != nil {
		return nil, err
	}

	b64sig := base64.URLEncoding.EncodeToString(*signature)

	jwt := tbs + "." + b64sig
	return &jwt, nil

}

func (s *Service) ValidateSession(jwt string) error {
	b64cert, err := s.kvStore.Get("cert")
	if err != nil {
		return err
	}
	rawCert, err := base64.StdEncoding.DecodeString(*b64cert)
	if err != nil {
		return err
	}
	certificate, err := x509.ParseCertificate(rawCert)
	if err != nil {
		return err
	}

	log.Println(certificate)

	// validate signature

	return nil
}

func encode(body []byte) string {
	return base64.URLEncoding.EncodeToString(body)
}

func sign(pkey *rsa.PrivateKey, tbs []byte) (*[]byte, error) {

	sig, err := rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA256, tbs)
	if err != nil {
		return nil, err
	}
	return &sig, nil

}
