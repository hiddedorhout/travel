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
	"errors"
	"fmt"
	"strings"
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
	Iat int    `json:"iat"`
	Exp int    `json:"exp"`
}

func (s *Service) GenerateSession(aud string) (*string, error) {

	b64pkey, err := s.kvStore.Get("pkey")
	if err != nil {
		return nil, errorHandler(err, "getting pkey")
	}

	rawPkey, err := base64.StdEncoding.DecodeString(*b64pkey)
	if err != nil {
		return nil, errorHandler(err, "decoding pkey")
	}
	pkey, err := x509.ParsePKCS1PrivateKey(rawPkey)
	if err != nil {
		return nil, errorHandler(err, "parsing pkey")
	}

	jwt, err := generateSession(aud, pkey)
	if err != nil {
		return nil, errorHandler(err, "creating jwt")
	}

	return jwt, nil

}

func (s *Service) ValidateSession(jwt string) error {
	b64cert, err := s.kvStore.Get("cert")
	if err != nil {
		return errorHandler(err, "Certificate")
	}
	rawCert, err := base64.StdEncoding.DecodeString(*b64cert)
	if err != nil {
		return errorHandler(err, "decoding cert")
	}
	certificate, err := x509.ParseCertificate(rawCert)
	if err != nil {
		return errorHandler(err, "parsing cert")
	}

	pubkey := certificate.PublicKey.(*rsa.PublicKey)

	if err := validateSession(jwt, pubkey); err != nil {
		return err
	}
	return nil
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

func validateSession(jwt string, pubkey *rsa.PublicKey) error {

	jot := strings.Split(jwt, ".")
	if len(jot) != 3 {
		return errors.New("Invalid jwt")
	}

	rawHeader, err := base64.URLEncoding.DecodeString(jot[0])
	if err != nil {
		return errorHandler(err, "base64 decode header")
	}

	rawPayload, err := base64.URLEncoding.DecodeString(jot[1])
	if err != nil {
		return errorHandler(err, "base64 decode payload")
	}

	signature, err := base64.URLEncoding.DecodeString(jot[2])
	if err != nil {
		return errorHandler(err, "base64 decode signature")
	}

	var header jwtHeader
	var payload jwtPayload

	if err := json.Unmarshal(rawHeader, &header); err != nil {
		return errorHandler(err, "Marshal header")
	}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return errorHandler(err, "Marshal payload")
	}

	tbs := jot[0] + "." + jot[1]

	h := sha256.New()
	h.Write([]byte(tbs))
	hashed := h.Sum(nil)

	if err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed, signature); err != nil {
		return errorHandler(err, "Invalid signature")
	}

	if payload.Exp > time.Now().Nanosecond() {
		return errorHandler(errors.New("Token error"), "token expired")
	}

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

func errorHandler(err error, message string) error {
	return fmt.Errorf("%s: %s", message, err.Error())
}
