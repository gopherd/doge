package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Payload holds the custom fields of jwt
type Payload struct {
	Salt   string   `json:"salt,omitempty"`
	Ver    string   `json:"ver,omitempty"`
	Uid    int64    `json:"uid,omitempty"`
	IP     string   `json:"ip,omitempty"`
	Loc    string   `json:"loc,omitempty"`
	Chan   int      `json:"chan,omitempty"`
	Os     string   `json:"os,omitempty"`
	Flags  int64    `json:"flags,omitempty"`
	Scopes []string `json:"scopes,omitempty"`

	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Gender int    `json:"gender,omitempty"`

	Accounts map[string]string `json:"accounts"`
}

// HasScope checks whether the payload has specified scope
func (p Payload) HasScope(scope string) bool {
	for _, s := range p.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// AddScopes appends a score
func (p *Payload) AddScopes(scopes ...string) {
	p.Scopes = append(p.Scopes, scopes...)
}

// Claims represents all the fields of jwt
type Claims struct {
	jwt.StandardClaims
	Payload
}

// Valid checks whether the claims is valid
func (c Claims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc().Unix()

	if c.VerifyExpiresAt(now, true) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyIssuedAt(now, false) == false {
		vErr.Inner = fmt.Errorf("token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

// Verifier ...
type Verifier struct {
	ecdsaPubkey *ecdsa.PublicKey
	keyId       string
}

func NewVerifier(filename, keyId string) (*Verifier, error) {
	key, err := loadPublicKey(filename)
	if err != nil {
		return nil, err
	}
	return &Verifier{
		ecdsaPubkey: key,
		keyId:       keyId,
	}, nil
}

func (v *Verifier) Verify(issuer, token string) (*Claims, error) {
	var claims = new(Claims)
	_, err := jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
		kid, ok := tok.Header["kid"]
		if !ok || kid == nil {
			return nil, jwt.ErrInvalidKey
		}
		if s, ok := kid.(string); !ok || s != v.keyId {
			return nil, jwt.ErrInvalidKey
		}
		return v.ecdsaPubkey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims.VerifyIssuer(issuer, true) == false {
		return nil, fmt.Errorf("jwt: token issuer mismatched")
	}
	return claims, nil
}

// Signer
type Signer struct {
	Verifier
	ecdsaKey *ecdsa.PrivateKey
}

func NewSigner(filename, keyId string) (*Signer, error) {
	key, err := loadPrivateKey(filename)
	if err != nil {
		return nil, err
	}
	return &Signer{
		Verifier: Verifier{
			ecdsaPubkey: &key.PublicKey,
			keyId:       keyId,
		},
		ecdsaKey: key,
	}, nil
}

// Sign signs the claims as a token string
func (s *Signer) Sign(claims *Claims) (string, error) {
	var tok = jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tok.Header["kid"] = s.keyId
	return tok.SignedString(s.ecdsaKey)
}

func loadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(bytes)
}

func parsePrivateKey(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("jwt: invalid ecdsa private key file")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pk := key.(type) {
	case *ecdsa.PrivateKey:
		return pk, nil
	default:
		return nil, errors.New("jwt: invalid ecdsa private key file")
	}
}

func loadPublicKey(filename string) (*ecdsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parsePublicKey(bytes)
}

func parsePublicKey(bytes []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("jwt: invalid ecdsa private key file")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pk := key.(type) {
	case *ecdsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("jwt: invalid ecdsa private key file")
	}
}
