package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) removePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	User_ID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	Post_ID, err := strconv.ParseUint(ps.ByName("postID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	Post, err := rt.db.GetPost(Post_ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting post by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h := r.Header.Get("Authorization")
	LoggedUserID, err := strconv.ParseUint(strings.TrimSpace(strings.Replace(h, "Bearer", "", 1)), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting token from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if User_ID != LoggedUserID {
		rt.baseLogger.Error("You can't delete someone else post")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = rt.db.RemovePost(User_ID, Post_ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error db removing post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = os.Remove(filepath.Join("/tmp/", Post.Media_file))
	if err != nil {
		ctx.Logger.WithError(err).Error("Post not found")
	}

	w.WriteHeader(http.StatusNoContent)
}
