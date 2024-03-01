package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error
	var user database.User
	var follower []database.User
	var following []database.User
	var posts []database.Post

	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting parameters id from string to int")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// check if the user requested exists
	userExist, err := rt.db.GetUserByID(id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting user by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userExist == (database.User{}) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h := r.Header.Get("Authorization")
	reqUserID, err := strconv.ParseUint(strings.TrimSpace(strings.Replace(h, "Bearer", "", 1)), 10, 64)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error converting token from string to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Requesting user is banned from target user
	ban, err := rt.db.IsBanned(reqUserID, id)
	if err != nil {
		ctx.Logger.WithError(err).Error("Banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user, err = rt.db.GetUserByID(id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting user by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	follower, err = rt.db.GetFollower(id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting follower by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	following, err = rt.db.GetFollowing(id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting following by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	posts, err = rt.db.GetPosts(id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting posts by ID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	ban, err = rt.db.IsBanned(id, reqUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error checking for ban")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send to the user
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(UserProfile{
		User_ID:   user.ID,
		Username:  user.Username,
		Followers: follower,
		Following: following,
		Posts:     posts,
		IsBanned:  ban,
	})
	if err != nil {
		rt.baseLogger.WithError(err).Warning("response return an error on encode")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
	}
}
