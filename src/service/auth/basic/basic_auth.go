package basic

import "net/http"

//Auth ... a holder around our functions.
type Auth struct {
}

//NewAuth ... creates a basic auth object
func NewAuth() *Auth {
	return &Auth{}
}

//Authorize ... handles a basic auth implementation.
func (a *Auth) Authorize(req *http.Request) (string, string, bool) {
	return req.BasicAuth()
}

//ValidateTokenHeader ...
func (a *Auth) ValidateTokenHeader(req *http.Request) (bool, error) {
	return false, nil
}

func (a *Auth) GenerateToken(input map[string]interface{}) (string, error) {
	return "", nil
}
