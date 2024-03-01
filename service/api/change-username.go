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

func (rt *_router) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error
	var user database.User

	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	h := r.Header.Get("Authorization")
	loggedUserID, err := strconv.ParseUint(strings.TrimSpace(strings.Replace(h, "Bearer", "", 1)), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting token from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// You can't change username for someone else
	if id != loggedUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	body := Username{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		rt.baseLogger.WithError(err).Warning("wrong JSON received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user.ID = id
	user.Username = strings.ToUpper(body.Username)

	actualUsername, err := rt.db.GetUser(body.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("error getting user by username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if actualUsername.ID != 0 {
		rt.baseLogger.Error("username already exists")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Username already exists, try another one!"})
		return
	}

	err = rt.db.ChangeUsername(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
