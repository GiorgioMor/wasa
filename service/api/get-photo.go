package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	postID, err := strconv.ParseUint(ps.ByName("postID"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error parsing postID from URL")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, err := rt.db.GetPost(postID)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error getting post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, "/tmp/"+post.Media_file)
}
