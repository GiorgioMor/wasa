package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Protinus/homework/service/api/reqcontext"
)

func (rt *_router) checkToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) bool {
	//	legge l'header e controlla se il token nell'header corrisponde all'id dell'utente
	h := r.Header.Get("Authorization")
	stringIDToken := strings.TrimSpace(strings.Replace(h, "Bearer", "", 1))
	if stringIDToken != "" {
		intIDToken, errS := strconv.ParseInt(stringIDToken, 10, 64)
		if errS != nil {
			rt.baseLogger.WithError(errS).Error("Error converting Token from string to int")
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		if intIDToken != -1 {
			user, errU := rt.db.GetUserByID(uint64(intIDToken))

			if user.Username == "" {
				rt.baseLogger.WithError(errU).Error("User un-authorized")
				w.WriteHeader(http.StatusUnauthorized)
				return false
			} else {
				return true
			}
		}
	}
	rt.baseLogger.Error("User un-authorized")
	w.WriteHeader(http.StatusUnauthorized)
	return false
}
