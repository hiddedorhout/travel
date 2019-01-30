package sessions

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"testing"
)

func TestGenerateSession(t *testing.T) {

	pkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	jwt, err := generateSession("Rachel", pkey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*jwt)
	t.Log("------------------")

	pubkey := pkey.PublicKey

	ppub := x509.MarshalPKCS1PublicKey(&pubkey)

	p := &pem.Block{
		Type:  "Public Key",
		Bytes: ppub,
	}
	if err := pem.Encode(os.Stdout, p); err != nil {
		log.Fatal(err)
	}

	if err := validateSession(*jwt, &pubkey); err != nil {
		t.Fatal(err)
	}
}
