package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) removeComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error

	PostsUserID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting URL parameters 'id' from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	commentID, err := strconv.ParseUint(ps.ByName("commentID"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting URL parameters 'commentID'  from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h := r.Header.Get("Authorization")
	UserID, err := strconv.ParseUint(strings.TrimSpace(strings.Replace(h, "Bearer", "", 1)), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting token from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The author of the post is removing any comment
	if PostsUserID == UserID {
		err = rt.db.RemoveComment(commentID)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error removing comment")
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Requesting user is banned from target user
	ban, err := rt.db.IsBanned(UserID, PostsUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	c, err := rt.db.GetCommentByID(commentID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting comment by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if c.User_ID == UserID {
		err = rt.db.RemoveComment(commentID)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error removing comment")
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		ctx.Logger.WithError(err).Error("Trying to remove someone else comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
