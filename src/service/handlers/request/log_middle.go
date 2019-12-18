package request

import (
	"fmt"
	"net/http"
	"service/log"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

var logClient log.ProdInterface

//Logger ... defines a response handler and waits for completion.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		t1 := time.Now()
		requestID := RetreiveRequestID(req.Context())

		if logClient != nil {
			logClient.Info("INCOMING", zap.String("requestID", requestID),
				zap.String("at", t1.String()))
		}

		defer func() {
			if logClient != nil {
				logClient.Info("RESPONSE",
					zap.Int("status", ww.Status()),
					zap.String("requestID", requestID),
					zap.String("size", fmt.Sprintf("%d bytes", ww.BytesWritten())),
					zap.String("in", time.Since(t1).String()),
				)
			}
		}()

		next.ServeHTTP(ww, req)
	})
}

//SetupLogger ... passes in injected clients
func SetupLogger(log log.ProdInterface) {
	logClient = log
}
