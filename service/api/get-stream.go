package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error
	var posts []database.Post

	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting parameters id from string to int")
		w.WriteHeader(http.StatusInternalServerError)
	}

	followers, err := rt.db.GetFollowing(id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting follower by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	for _, f := range followers {

		ps, err := rt.db.GetPosts(f.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = append(posts, ps...)

	}

	// Send the stream to the user.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	errr := json.NewEncoder(w).Encode(UserHome{
		Followers: followers,
		Posts:     posts,
	})
	if errr != nil {
		rt.baseLogger.WithError(errr).Warning("posts retun an error on encode")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
	}
}
