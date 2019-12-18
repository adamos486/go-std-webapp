package database

import "database/sql"

//DBInterface declares an interface that adheres with the sql lib definition
//go:generate counterfeiter . DBInterface
type DBInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

//Client defines an object that binds methods using a passed in DBInterface
type Client struct {
	DB DBInterface
}

//New creates a new DB with a bound passed in DBInterface
func New(db DBInterface) *Client {
	return &Client{DB: db}
}
