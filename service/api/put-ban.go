package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error

	Ban_ID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// check if the user requested exists
	userExist, err := rt.db.GetUserByID(Ban_ID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting user by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userExist == (database.User{}) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	User_ID, err := strconv.ParseUint(ps.ByName("bUserID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// You can't ban yourselft
	if Ban_ID == User_ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.BanUser(Ban_ID, User_ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Unfollow the user banned and remove following too
	err = rt.db.UnfollowUser(Ban_ID, User_ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = rt.db.UnfollowUser(User_ID, Ban_ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
