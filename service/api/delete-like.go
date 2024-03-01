package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) unlikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error

	PostsUserID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting URL parameters 'id' from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	PostID, err := strconv.ParseUint(ps.ByName("postID"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting URL parameters 'postID'  from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	likeID, err := strconv.ParseUint(ps.ByName("likeID"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting URL parameters 'likeID'  from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if PostsUserID == likeID {
		ctx.Logger.WithError(err).Error("trying to unlike your own post")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Requesting user is banned from target user
	ban, err := rt.db.IsBanned(likeID, PostsUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	u, err := rt.db.GetUserByID(likeID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting user by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = rt.db.UnlikePost(PostID, u)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error removing like")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
