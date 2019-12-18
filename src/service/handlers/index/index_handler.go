package index

import (
	"encoding/json"
	"io"
	"net/http"
	"service/database"
	"service/log"
	"time"

	"go.uber.org/zap"
)

//Handler ... contains all handlers for index route.
//go:generate counterfeiter . Handler
type Handler interface {
	Handler(w http.ResponseWriter, req *http.Request)
}

//EventsResponse ... is the current representation of the response json.
type EventsResponse struct {
	Code int        `json:"code"`
	List []EventRow `json:"list"`
}

//EventRow ... is the current database representation.
type EventRow struct {
	ID          int       `pq:"id" json:"id"`
	Name        string    `pq:"name" json:"name"`
	Description string    `pq:"description" json:"description"`
	DateAdded   time.Time `pq:"date_added" json:"dateAdded"`
}

//Index ... holds a logger, a dbClient, and an auth service.
type Index struct {
	log      log.ProdInterface
	dbClient database.DBInterface
}

//Handler is the handler for the root path.
func (i *Index) Handler(w http.ResponseWriter, req *http.Request) {
	i.veryBadNesting(
		i.passThrough(func(w http.ResponseWriter, req *http.Request) {
			i.indexLogic(w, req)
		}),
	)(w, req)
}

//New ... returns a pointer to a new Index object.
func New(log log.ProdInterface, db database.DBInterface) *Index {
	return &Index{
		log:      log,
		dbClient: db,
	}
}

func (i *Index) indexLogic(w http.ResponseWriter, req *http.Request) {
	rows, err := i.dbClient.Query("SELECT id, name, description, date_added FROM event;")
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()
	eventRows := make([]EventRow, 0)
	for rows.Next() {
		var eventRow EventRow
		if scanErr := rows.Scan(
			&eventRow.ID,
			&eventRow.Name,
			&eventRow.Description,
			&eventRow.DateAdded); scanErr != nil {
			panic(scanErr)
		}
		eventRows = append(eventRows, eventRow)
	}

	jsonErr := marshalEventRows(eventRows, w)
	if jsonErr != nil {
		panic(jsonErr)
	}
}

func marshalEventRows(events []EventRow, w io.Writer) error {
	response := EventsResponse{
		Code: http.StatusOK,
		List: events,
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func (i *Index) passThrough(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id, ok := ctx.Value(42).(int64)
		if ok {
			i.log.Debug("passThrough::id->", zap.Int64("id", id))
		}
		i.log.Debug("PassThrough")
		next(w, req)
	}
}

func (i *Index) veryBadNesting(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		i.log.Debug("Very bad case of nesting")
		next(w, req)
	}
}
