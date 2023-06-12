package handlers

import (
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

// UsersRouter handles requests for the users route
func UsersRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			usersGetAll(w, r)
			return
		case http.MethodHead:
			usersGetAll(w, r)
			return
		case http.MethodPost:
			usersPostOne(w, r)
			return
		case http.MethodOptions:
			postOptionsResponse(w, []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodOptions}, nil)
			return
		default:
			postError(w, http.StatusMethodNotAllowed)
		}
	}
	path = strings.TrimPrefix(path, "/users/")
	if !bson.IsObjectIdHex(path) {
		postError(w, http.StatusNotFound)
		return
	}

	id := bson.ObjectIdHex(path)

	switch r.Method {
	case http.MethodGet:
		usersGetOne(w, r, id)
		return
	case http.MethodHead:
		usersGetOne(w, r, id)
		return
	case http.MethodPut:
		usersPutOne(w, r, id)
		return
	case http.MethodPatch:
		usersPatchOne(w, r, id)
		return
	case http.MethodDelete:
		usersDeleteOne(w, r, id)
		return
	case http.MethodOptions:
		postOptionsResponse(w, []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, nil)
		return
	default:
		postError(w, http.StatusMethodNotAllowed)
	}
}
