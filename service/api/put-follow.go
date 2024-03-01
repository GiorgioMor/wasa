package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	User_ID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// check if the user requested exists
	userExist, err := rt.db.GetUserByID(User_ID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting user by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userExist == (database.User{}) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	Follow_ID, err := strconv.ParseUint(ps.ByName("fUserID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Requesting user is banned from target user
	ban, err := rt.db.IsBanned(Follow_ID, User_ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// You can't follow yourself
	if User_ID == Follow_ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errr := rt.db.FollowUser(User_ID, Follow_ID)
	if errr != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
