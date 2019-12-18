package identity

import (
	"encoding/json"
	"io"
	"time"
)

//Row ... this is a data struct to house a database response.
type Row struct {
	ID          string      `pq:"id" json:"id"`
	FirstName   string      `pq:"first_name" json:"firstName"`
	LastName    string      `pq:"last_name" json:"lastName"`
	ProfileInfo interface{} `pq:"profile" json:"profile"`
	CreatedAt   time.Time   `pq:"created_at" json:"createdAt"`
	UpdatedAt   time.Time   `pq:"updated_at" json:"updatedAt"`
}

type singularResponse struct {
	Code    int `json:"status"`
	Element Row `json:"identity"`
}

//RespondWithJSON ... marshals this Row into JSON and returns it.
func (i *Row) RespondWithJSON(status int, w io.Writer) error {
	response := singularResponse{
		Code:    status,
		Element: *i,
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}
