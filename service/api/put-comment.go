package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) commentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	h := r.Header.Get("Authorization")
	commentUserID, err := strconv.ParseUint(strings.TrimSpace(strings.Replace(h, "Bearer", "", 1)), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting token from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Requesting user is banned from target user
	ban, err := rt.db.IsBanned(commentUserID, PostsUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	u, err := rt.db.GetUserByID(commentUserID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting user by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	body := CommentText{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		rt.baseLogger.WithError(err).Warning("wrong JSON received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = rt.db.AddComment(PostID, u, body.Comment)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error adding comment")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
