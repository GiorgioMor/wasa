package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

// login
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var err error
	var user database.User

	body := Username{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		rt.baseLogger.WithError(err).Warning("wrong JSON received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	username := strings.ToUpper(body.Username)

	user, err = rt.db.GetUser(username)

	if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == (database.User{}) {
		id, err := rt.db.CreateUser(username)

		if err != nil {
			ctx.Logger.WithError(err).Error("can't create user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			errr := json.NewEncoder(w).Encode(id)
			if errr != nil {
				rt.baseLogger.WithError(errr).Warning("id (1) return an error on encode")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
			}
		}
	} else {
		id := user.ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(id)
		if err != nil {
			rt.baseLogger.WithError(err).Warning("id (2) retursn an error on encode")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		}
	}
}
