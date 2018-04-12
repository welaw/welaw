package oauth

import (
	"crypto/subtle"
	"time"

	stdjwt "github.com/dgrijalva/jwt-go"
)

const (
	kid = "kid"
)

type CustomClaims struct {
	//Kid      string `json:"kid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	stdjwt.StandardClaims
}

func (c CustomClaims) Valid() error {
	return nil
}

func (c CustomClaims) VerifyUsername(u string) bool {
	return subtle.ConstantTimeCompare([]byte(c.Username), []byte(u)) != 0
}

var (
	method         = stdjwt.SigningMethodHS256
	testSigningKey = []byte("test_signing_key")
	testClaims     = CustomClaims{}
)

func MakeToken(_, uid, name, email string, key []byte) (string, error) {

	claims := CustomClaims{
		Username: name,
		//Kid:      kid,
		StandardClaims: stdjwt.StandardClaims{
			Subject:   uid,
			Audience:  "welaw",
			Issuer:    "welaw.auth",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)

	jwtString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func MakeClaimsToken(uid string) *stdjwt.Token {
	claims := stdjwt.StandardClaims{
		Subject:   uid,
		Audience:  "welaw",
		Issuer:    "welaw.auth",
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	return stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)
}

func MakeRefreshToken(_, uid, role string, key []byte) (string, error) {

	claims := CustomClaims{
		Username: uid,
		Role:     role,
		//Kid:      kid,
		StandardClaims: stdjwt.StandardClaims{
			Subject:   uid,
			Audience:  "welaw",
			Issuer:    "welaw.auth",
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
	}

	token := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)

	jwtString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
