package basic

import (
	"net/http"
	"service/auth"
	"service/log"
)

//AuthMiddleware ... performs basic auth
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		u, p, hasAuth := authClient.Authorize(req)
		if hasAuth && u == "tony" && p == "house" {
			next.ServeHTTP(w, req)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}

var authClient *auth.Client
var logClient log.ProdInterface

//SetupAuthMiddleware ... attaches a configured authClient
func SetupAuthMiddleware(auth *auth.Client, log log.ProdInterface) {
	authClient = auth
	logClient = log
}
