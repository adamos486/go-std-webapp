package loggederror

import (
	"net/http"
	"service/handlers/request"
	"service/log"

	"go.uber.org/zap"
)

//RespondWithProperErrorAndLogIt ...
//will respond with an error object that is marshaled to json, and wrap the message from
//the passed in error.
func RespondWithProperErrorAndLogIt(log log.ProdInterface, status int,
	err error, context string, w http.ResponseWriter, req *http.Request) {
	//proper http error using standard lib.
	requestID := request.RetreiveRequestID(req.Context())
	if err != nil {
		log.Error(context, zap.String("requestID", requestID), zap.Error(err))
		http.Error(w, err.Error(), status)
		return
	}
	log.Info(context, zap.String("requestID", requestID), zap.String("error",
		"GENERIC ERROR!!!"))
	http.Error(w, "generic error", status)
}

// RespondWithWithExpectedSoftError ...
// will respond to the http.Request with a public message, while logging the
// requestID, message, and source.
func RespondWithWithExpectedSoftError(log log.ProdInterface, status int, message string,
	source string, w http.ResponseWriter, req *http.Request) {
	requestID := request.RetreiveRequestID(req.Context())
	log.Info("Logic Error->", zap.String("requestID", requestID),
		zap.String("message", message), zap.String("source", source))
	http.Error(w, message, status)
}
