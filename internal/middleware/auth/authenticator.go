package auth

import "crypto/subtle"

type Authenticator interface {
	Authenticate(token string) bool
}

type TokenComparer func([]byte, []byte) bool

type TokenAuthenticator struct {
	expectedToken string
}

func NewTokenAuthenticator(token string) *TokenAuthenticator {
	return &TokenAuthenticator{expectedToken: token}
}

func (a *TokenAuthenticator) Authenticate(token string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(a.expectedToken)) == 1
}
