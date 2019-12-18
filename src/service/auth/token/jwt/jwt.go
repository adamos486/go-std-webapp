package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
)

//Interface ...
//JWT Interface for contract testing.
type Interface interface {
	ValidateToken(token string) (map[string]interface{}, bool, error)
	Generate(input map[string]interface{}) (string, error)
}

//Service ...
//A holder struct for a jwt impl.
type Service struct {
}

//NewService ...
//Creates a new jwt token service.
func NewService() *Service {
	return &Service{}
}

//IdentityClaims ...
//Holds personalized identity information as well as the standard information.
type IdentityClaims struct {
	Email string `json:"email"`
	jwtGo.StandardClaims
}

//ValidateToken ...
//Implements method on TokenService interface with real JWT logic.
func (s *Service) ValidateToken(token string) (interface{}, bool, error) {
	if token == "" {
		return nil, false, errors.New("cannot validate an empty token")
	}
	tokenObj, err := jwtGo.ParseWithClaims(token, &IdentityClaims{}, obtainJwtKey)
	if err != nil {
		return nil, false, err
	}
	return tokenObj, tokenObj.Valid, nil
}

//Generate ...
//This generates a jwt token with the passed in values.
func (s *Service) Generate(input map[string]interface{}) (string, error) {
	if len(input) == 0 {
		return "", errors.New("can't generate a token with required claims")
	}
	claims := IdentityClaims{}
	if val, ok := input["email"]; ok {
		claims.Email = val.(string)
	}
	claims.ExpiresAt = time.Now().Add(12 * time.Hour).UTC().Unix()
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	if token != nil {
		tokenString, err := token.SignedString([]byte(os.Getenv("STAGE_JWT_SECRET")))
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}
	return "", nil
}

func obtainJwtKey(token *jwtGo.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("STAGE_JWT_SECRET")), nil
}
