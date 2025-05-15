package auth

type Authenticator interface {
	Authenticate(token string) bool
}

type TokenAuthenticator struct {
	expectedToken string
}

func NewTokenAuthenticator(token string) *TokenAuthenticator {
	return &TokenAuthenticator{expectedToken: token}
}

func (a *TokenAuthenticator) Authenticate(token string) bool {
	return token == a.expectedToken
}
