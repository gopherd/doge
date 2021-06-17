package jwt_test

import (
	"reflect"
	"testing"
	"time"

	. "github.com/gopherd/doge/jwt"
)

func TestSignerVerifier(t *testing.T) {
	const (
		keyId  = "testing"
		issuer = "gopherd.com"
	)

	signer, err := NewSigner("testdata/ec256.p8", keyId)
	if err != nil {
		t.Fatalf("NewSigner error: %v", err)
	}
	verifier, err := NewVerifier("testdata/ec256.pub.p8", keyId)
	if err != nil {
		t.Fatalf("NewVerifier error: %v", err)
	}

	// claims data
	claims := &Claims{
		Payload: Payload{
			Salt: "abc",
			Ver:  "1",
			Uid:  124,
			Name: "hello",
		},
	}
	claims.Issuer = issuer
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = claims.IssuedAt + 3600 // 1h

	// Sign
	token, err := signer.Sign(claims)
	if err != nil {
		t.Fatalf("signer.Sign error: %v", err)
	}

	// Verify
	got, err := verifier.Verify(issuer, token)
	if err != nil {
		t.Fatalf("signer.Sign error: %v", err)
	}

	// assert
	if !reflect.DeepEqual(claims, got) {
		t.Errorf("want %v, got %v", claims, got)
	}
}
