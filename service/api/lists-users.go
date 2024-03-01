package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

func (rt *_router) listUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	rt.checkToken(w, r, ps, ctx)

	var err error
	var users []database.User

	// Extract the query parameter from the URL
	username := r.URL.Query().Get("username")

	h := r.Header.Get("Authorization")
	stringIDToken := strings.TrimSpace(strings.Replace(h, "Bearer", "", 1))
	uintIDToken, _ := strconv.ParseUint(stringIDToken, 10, 64)

	if username != "" {
		// Search the user in the database filtered by the query parameter
		users, err = rt.db.SearchUser(username, uintIDToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Database has encountered an error")
			_ = json.NewEncoder(w).Encode([]User{})
			return
		}
	} else {
		// Request an unfiltered list of users from the DB
		users, err = rt.db.SearchUser("", uintIDToken)
		if err != nil {
			// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
			// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
			ctx.Logger.WithError(err).Error("can't list users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Send the list to the user.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	errr := json.NewEncoder(w).Encode(users)
	if errr != nil {
		rt.baseLogger.WithError(errr).Warning("users retun an error on encode")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
	}
}
