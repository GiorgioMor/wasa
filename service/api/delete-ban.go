package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	Ban_ID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	User_ID, err := strconv.ParseUint(ps.ByName("bUserID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// You can't unban yourself
	if Ban_ID == User_ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errr := rt.db.UnbanUser(Ban_ID, User_ID)
	if errr != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
