package fail

import (
	"nbientry/web/common"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
)

const (
	PUBLIC = iota
	VIEWER
	ADMIN
)

var ValidIDs = regexp.MustCompile(`^([0-9a-z\_\-]+)$`)

func Auth(w http.ResponseWriter, r *http.Request, user common.User, level int) bool {
	if level >= VIEWER && user.Email == "" {
		return throw_auth(w, r, "Not logged in")
	}

	if level >= ADMIN && !user.Admin {
		return throw_auth(w, r, "Not an admin")
	}

	return false
}

func throw_auth(w http.ResponseWriter, r *http.Request, reason string) bool {
	w.Header().Set("X-Auth-Reason", reason)
	http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
	return true
}

func Render(w http.ResponseWriter, r *http.Request, elem templ.Component) {
	err := elem.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Form(w http.ResponseWriter, r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return true
	}

	return false
}
