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
}

func TestOther(t *testing.T) {
	s := "eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ"
	t.Log([]byte(s))
}
