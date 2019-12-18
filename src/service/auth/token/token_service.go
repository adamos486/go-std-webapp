package token

//Interface ... a interface definition for implementers
//go:generate counterfeiter . Interface
type Interface interface {
	ValidateToken(token string) (interface{}, bool, error)
	Generate(input map[string]interface{}) (string, error)
}

//Service ...
//This is a generic token service that can take in any type of adhering typed token.
type Service struct {
	I Interface
}

//NewService ...
//Creates a new instance of Token Service and passes back a pointer to it.
func NewService(typedToken Interface) *Service {
	return &Service{
		I: typedToken,
	}
}

//ValidateToken ...
//Adhering method for validating a generic token -> typed token
func (s *Service) ValidateToken(token string) (interface{}, bool, error) {
	return s.I.ValidateToken(token)
}

//Generate ...
//Generates a new JWT token and returns it, otherwise returns an error.
func (s *Service) Generate(input map[string]interface{}) (string, error) {
	return s.I.Generate(input)
}
