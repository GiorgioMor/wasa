package api

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
	"gitlab.com/Protinus/homework/service/database"
)

func (rt *_router) newPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.checkToken(w, r, ps, ctx)

	var err error
	var post database.Post

	userID, err := strconv.ParseUint(ps.ByName("id"), 10, 64)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("fail convert id from string to uint")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error Parse Multipart Form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error formfile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("/tmp/"+ps.ByName("id")+"_"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error creating file")
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving image")
		return
	}

	post.Caption = r.FormValue("caption")
	post.Created_Datetime = time.Now().Format("2006-01-02 15:04:05")
	post.User_ID = userID
	post.Media_file = ps.ByName("id") + "_" + r.FormValue("fileName")

	err = rt.db.CreatePost(post)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't create post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
