package identity

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"service/auth"
	"service/handlers/loggederror"
	"service/log"

	"go.uber.org/zap"
)

//HandlerInterface ... contains all handlers for identity route
//go:generate counterfeiter . HandlerInterface
type HandlerInterface interface {
	Handler(w http.ResponseWriter, req *http.Request)
	CreateIdentity(w http.ResponseWriter, req *http.Request)
	AuthIdentity(w http.ResponseWriter, req *http.Request)
}

//HandlerObject ... holds elementals for interface methods.
type HandlerObject struct {
	Log     log.ProdInterface
	Service ServiceInterface
	Auth    auth.Interface
}

//NewHandlerObject ... returns a pointer to a new Identity Object
func NewHandlerObject(logClient log.ProdInterface, service ServiceInterface,
	auth auth.Interface) *HandlerObject {
	return &HandlerObject{
		Log:     logClient,
		Service: service,
		Auth:    auth,
	}
}

//CreateIdentity ...
//This is the router handler Func for identity creation.
//This handler will validate input, validate auth, and if approved, pass
//info from the POST body to the identity_service
func (h *HandlerObject) CreateIdentity(w http.ResponseWriter, req *http.Request) {
	identity, result, err := h.Service.Create("uuidv7")
	if err != nil {
		loggederror.RespondWithProperErrorAndLogIt(
			h.Log,
			http.StatusInternalServerError,
			err,
			"identity_handler::CreateIdentity",
			w,
			req,
		)
		return
	}
	if identity != nil {
		h.Log.Debug("CreateIdentity", zap.Any("identity", identity),
			zap.Any("result", result))
		responseErr := identity.RespondWithJSON(http.StatusOK, w)
		if responseErr != nil {
			h.Log.Error("CreateIdentity::RespondWithJSON", zap.Error(responseErr))
		}
		return
	}
	loggederror.RespondWithProperErrorAndLogIt(
		h.Log,
		http.StatusInternalServerError,
		err,
		"identity_handler::Service created nil identity",
		w,
		req,
	)
}

type authPostBody struct {
	ID                string `json:"id"`
	EncryptedPassword string `json:"password"`
}

// internalServerError is used to wrap our loggederror for this route.
func (h *HandlerObject) internalServerError(err error, source string, w http.ResponseWriter,
	req *http.Request) {
	loggederror.RespondWithProperErrorAndLogIt(
		h.Log,
		http.StatusInternalServerError,
		err,
		"identity_handler::"+source,
		w,
		req,
	)
}

// badRequest is used to wrap our loggederror for this route.
func (h *HandlerObject) badRequest(message, source string, w http.ResponseWriter,
	req *http.Request) {
	loggederror.RespondWithWithExpectedSoftError(
		h.Log,
		http.StatusBadRequest,
		message,
		"identity_handler::"+source,
		w,
		req,
	)
}

type authResponse struct {
	Status int    `json:"status"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

//AuthIdentity generates a jwt token for a known identity.
func (h *HandlerObject) AuthIdentity(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.internalServerError(err, "AuthIdentity", w, req)
		return
	}
	var jsonDoc authPostBody
	jsonErr := json.Unmarshal(body, &jsonDoc)
	if jsonErr != nil {
		h.internalServerError(jsonErr, "AuthIdentity", w, req)
		return
	}
	if jsonDoc.ID == "" || jsonDoc.EncryptedPassword == "" {
		h.badRequest("missing required auth params", "AuthIdentity", w, req)
		return
	}
	defer func() {
		closeErr := req.Body.Close()
		if closeErr != nil {
			h.internalServerError(closeErr, "AuthIdentity", w, req)
			return
		}
	}()
	input := make(map[string]interface{}, 0)
	input["ID"] = jsonDoc.ID
	token, tokenErr := h.Auth.GenerateToken(input)
	if tokenErr != nil {
		h.internalServerError(tokenErr, "AuthIdentity", w, req)
		return
	}

	response := authResponse{
		Status: http.StatusOK,
		ID:     jsonDoc.ID,
		Token:  token,
	}
	bytesArray, marshalErr := json.Marshal(&response)
	if marshalErr != nil {
		h.internalServerError(marshalErr, "AuthIdentity", w, req)
		return
	}

	_, writeErr := w.Write(bytesArray)
	if writeErr != nil {
		h.internalServerError(writeErr, "AuthIdentity", w, req)
	}
}

//Handler ... contains a handler funcction to be passed to our router.
func (h *HandlerObject) Handler(w http.ResponseWriter, req *http.Request) {
	row, err := h.Service.Fetch("uuidv7")
	if err != nil {
		loggederror.RespondWithProperErrorAndLogIt(h.Log, http.StatusInternalServerError,
			err, "identity_handler::fetch from service", w, req)
		return
	}
	h.Log.Debug("identity Handler", zap.Any("row", row))
	//h.Log.Info("we got a row!", zap.Any("row", row))
	if row != nil {
		err = row.RespondWithJSON(http.StatusOK, w)
		if err != nil {
			loggederror.RespondWithProperErrorAndLogIt(
				h.Log,
				http.StatusInternalServerError,
				err,
				"identity_handler::Handler",
				w,
				req)
			return
		}
	}
}
