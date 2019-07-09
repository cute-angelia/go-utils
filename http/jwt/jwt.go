package jwt

import (
	"github.com/robbert229/jwt"
)

func NewJwt(secret string) *Jwt {
	return &Jwt{
		Secret: secret,
	}
}

type Jwt struct {
	Secret string
}

func (self *Jwt) Generate(claims *jwt.Claims) (string, error) {
	algorithm := jwt.HmacSha256(self.Secret)
	token, err := algorithm.Encode(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (self *Jwt) Validate(token string) error {
	algorithm := jwt.HmacSha256(self.Secret)
	if err := algorithm.Validate(token); err != nil {
		return err
	} else {
		return nil
	}
}

func (self *Jwt) Decode(token string) (*jwt.Claims, error) {
	algorithm := jwt.HmacSha256(self.Secret)
	loadedClaims, err := algorithm.Decode(token)
	if err != nil {
		return nil, err
	} else {
		return loadedClaims, nil
	}
}
