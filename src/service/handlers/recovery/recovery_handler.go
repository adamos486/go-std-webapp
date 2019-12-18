package recovery

import (
	"net/http"
	"runtime/debug"
	"service/handlers/request"
	"service/log"

	"go.uber.org/zap"
)

var logClient log.ProdInterface

//Recover recovers from panics and logs the error with our proper logger.
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID := request.RetreiveRequestID(r.Context())
				logClient.Error("recovered from error", zap.ByteString("stack", debug.Stack()),
					zap.String("requestID", requestID))
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//SetupRecover passes in the log interface to establish the logger used.
func SetupRecover(log log.ProdInterface) {
	logClient = log
}
