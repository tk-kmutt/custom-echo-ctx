package jwt

import (
	"time"

	"golang.org/x/xerrors"

	"github.com/golang-jwt/jwt"
)

const (
	iat string = "iat"
	exp        = "exp"
	iss        = "iss"
	lu         = "lu"
)

type JWT struct {
	secret   []byte
	lifeTime int
	issuer   string
}

func NewJWT(secret []byte, lifeTime int, issuer string) *JWT {
	return &JWT{secret: secret, lifeTime: lifeTime, issuer: issuer}
}

type LoginUser struct {
	Name  string `json:"n"`
	Email string `json:"e"`
}

func (r *JWT) SignedString(user LoginUser) (string, error) {
	now := time.Now()
	// Set custom claims
	claims := jwt.MapClaims{
		iat: now.Unix(),
		exp: now.Add(time.Duration(r.lifeTime) * time.Hour).Unix(),
		iss: r.issuer,
		lu:  user,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString(r.secret)
}

func (r *JWT) Verify(signedString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", xerrors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return r.secret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, xerrors.Errorf("%s is expired", signedString)
			} else {
				return nil, xerrors.Errorf("%s is invalid", signedString)
			}
		} else {
			return nil, xerrors.Errorf("%s is invalid", signedString)
		}
	}

	if token == nil {
		return nil, xerrors.Errorf("not found token in %s:", signedString)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, xerrors.Errorf("not found claims in %s", signedString)
	}
	if mapClaims[iss] != r.issuer {
		return nil, xerrors.Errorf("unknown issuer %s in %s", mapClaims[iss], signedString)
	}

	user := mapClaims[lu]
	if user == nil {
		return nil, xerrors.Errorf("not found user in %s", signedString)
	}

	return mapClaims, nil
}

func (r *JWT) VerifySignedString(signedString string) (string, error) {
	mapClaims, err := r.Verify(signedString)
	if err != nil {
		return "", err
	}

	user := mapClaims[lu].(map[string]interface{})
	return r.SignedString(LoginUser{
		Name:  user["n"].(string),
		Email: user["e"].(string),
	})
}
