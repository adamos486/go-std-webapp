package request

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

type key string

const requestIDKey = key("requestID")

//GenerateRequestIDMiddle ... generates a random UUID + timestamp
//for each incoming request.
func GenerateRequestIDMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rightNow := time.Now()
		rightNowNano := rightNow.UnixNano()
		u1 := uuid.NewV4()
		requestID := fmt.Sprintf("%v-%d", u1.String(), rightNowNano)

		ctx := req.Context()
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

//RetreiveRequestID ... retrieves a requestID with private key from
//context and returns as a string.
func RetreiveRequestID(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}
