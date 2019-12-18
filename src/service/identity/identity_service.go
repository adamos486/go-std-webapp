package identity

import (
	"database/sql"
	"encoding/json"
	"errors"
	"service/database"
	"service/log"
	"time"

	"github.com/google/uuid"
	uuid2 "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

//ServiceInterface ... defines a required interface for all identity service methods.
//go:generate counterfeiter . ServiceInterface
type ServiceInterface interface {
	Fetch(id string) (*Row, error)
	Create(id string) (*Row, sql.Result, error)
}

//ServiceObject ...
//contains all elementals needing to be injected in tests, and used to perform business.
type ServiceObject struct {
	log log.ProdInterface
	db  database.DBInterface
}

//NewServiceObject ...
//takes in a logClient, dbClient and returns a pointer to a new ServiceObject that
//can contain fakes or real clients and perform bound methods.
func NewServiceObject(logClient log.ProdInterface, dbClient database.DBInterface) *ServiceObject {
	return &ServiceObject{
		log: logClient,
		db:  dbClient,
	}
}

type jsonObject struct {
	Email string `json:"email"`
}

//Create ...
//Creates a unique uuid.v4, creates a identity record, queries for the created record, and
//returns the created record.
func (s *ServiceObject) Create(id string) (*Row, sql.Result, error) {
	var identity Row
	jsonObj := jsonObject{
		Email: "test@gmail.com",
	}
	generatedID := uuid.New()
	generatedVariant := uuid2.NewV4()
	supraID := generatedVariant.String() + "-" + generatedID.String()
	supraID = supraID[0:50]
	s.log.Debug("supraID:", zap.String("supraID", supraID), zap.Int("supra lenght", len(supraID)))
	rawJSON, _ := json.Marshal(&jsonObj)
	rightNow := time.Now()
	result, err := s.db.Exec(`INSERT INTO identity
		(id, first_name, last_name, profile, created_at, updated_at) VALUES
		($1, $2, $3, $4, $5, $6);`,
		supraID,
		"adam",
		"cobb",
		rawJSON,
		rightNow,
		rightNow)
	if err != nil {
		return nil, result, err
	}

	//Examine the result and define errors if necessary.
	affected, resErr := result.RowsAffected()
	if resErr != nil {
		return nil, nil, resErr
	}
	if affected == int64(0) || affected > 1 {
		insertErr := errors.New("INSERT INTO had fatal errors")
		return nil, nil, insertErr
	}
	sqlRow := s.db.QueryRow(`SELECT id, first_name, last_name, profile, created_at, updated_at
		FROM identity WHERE id = $1`, supraID)
	if sqlRow != nil {
		var jsonData []byte

		scanErr := sqlRow.Scan(&identity.ID, &identity.FirstName, &identity.LastName,
			&jsonData, &identity.CreatedAt, &identity.UpdatedAt)
		if scanErr != nil {
			return nil, nil, scanErr
		}

		if len(jsonData) > 0 {
			var output interface{}
			decodeErr := json.Unmarshal(jsonData, &output)
			if decodeErr != nil {
				return nil, nil, decodeErr
			}
			identity.ProfileInfo = output
		}
	}
	return &identity, result, nil
}

//Fetch ... is an interface method for fetching identity records.
func (s *ServiceObject) Fetch(id string) (*Row, error) {
	row := s.db.QueryRow(
		"SELECT id, first_name, last_name, profile, created_at, updated_at FROM identity WHERE id = $1;", id)
	s.log.Debug("Fetch", zap.Any("row", row))
	if row != nil {
		var identityRow Row
		var jsonData []byte
		if err := row.Scan(&identityRow.ID, &identityRow.FirstName, &identityRow.LastName,
			&jsonData, &identityRow.CreatedAt, &identityRow.CreatedAt); err != nil {
			return nil, err
		}
		if len(jsonData) > 0 {
			var output interface{}
			decodingError := json.Unmarshal(jsonData, &output)
			if decodingError != nil {
				return nil, decodingError
			}
			identityRow.ProfileInfo = output
		}
		return &identityRow, nil
	}
	return nil, nil
}
