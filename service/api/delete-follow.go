package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	User_ID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	Follow_ID, err := strconv.ParseUint(ps.ByName("fUserID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// You can't unfollow yourself
	if User_ID == Follow_ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errr := rt.db.UnfollowUser(User_ID, Follow_ID)
	if errr != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
