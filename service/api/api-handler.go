package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	// Login
	rt.router.POST("/session", rt.wrap(rt.login))

	rt.router.GET("/users", rt.wrap(rt.listUsers))
	rt.router.DELETE("/users/:id", rt.wrap(rt.deleteUser))
	rt.router.PUT("/users/:id", rt.wrap(rt.changeUsername))
	rt.router.GET("/users/:id", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:id/home", rt.wrap(rt.getMyStream))

	// Post
	rt.router.POST("/users/:id/photos", rt.wrap(rt.newPost))
	rt.router.DELETE("/users/:id/photos/:postID", rt.wrap(rt.removePost))
	rt.router.GET("/users/:id/photos/:postID", rt.wrap(rt.getPhoto))

	// Comment/Uncomment
	rt.router.POST("/users/:id/photos/:postID/comments", rt.wrap(rt.commentPost))
	rt.router.DELETE("/users/:id/photos/:postID/comments/:commentID", rt.wrap(rt.removeComment))

	// Like/Unlike
	rt.router.PUT("/users/:id/photos/:postID/like/:likeID", rt.wrap(rt.likePost))
	rt.router.DELETE("/users/:id/photos/:postID/like/:likeID", rt.wrap(rt.unlikePost))

	// Follow/Unfollow
	rt.router.PUT("/users/:id/follow/:fUserID", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:id/follow/:fUserID", rt.wrap(rt.unfollowUser))

	// Ban/Unban
	rt.router.PUT("/users/:id/ban/:bUserID", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:id/ban/:bUserID", rt.wrap(rt.unbanUser))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
