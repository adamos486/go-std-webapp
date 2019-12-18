package auth

import (
	"net/http"
	"service/auth/token"
)

//Interface ... a interface definition for implementers
//go:generate counterfeiter . Interface
type Interface interface {
	Authorize(req *http.Request) (string, string, bool)
	ValidateTokenHeader(req *http.Request) (bool, error)
	GenerateToken(map[string]interface{}) (string, error)
}

//Client ... a client wrapper for basic and oauth.
type Client struct {
	I Interface
	T token.Interface
}

//NewClient ... creates a new Auth DB and returns a pointer to it.
func NewClient(auth Interface, token token.Interface) *Client {
	return &Client{
		I: auth,
		T: token,
	}
}

//Authorize ... performs and authorization and returns whether it's authorized.
func (c *Client) Authorize(req *http.Request) (string, string, bool) {
	return c.I.Authorize(req)
}

// GenerateToken ...
// generates a new JWT token with a given input map.
func (c *Client) GenerateToken(input map[string]interface{}) (string, error) {
	return c.T.Generate(input)
}

//ValidateTokenHeader ...
//Takes in a request, pulls of a header value called "token", passes it through
//jwt validation, returns the decoded object and/or whether or not the token is valid.
func (c *Client) ValidateTokenHeader(req *http.Request) (bool, error) {
	token := req.Header.Get("token")
	_, isValid, err := c.T.ValidateToken(token)
	return isValid, err
}
