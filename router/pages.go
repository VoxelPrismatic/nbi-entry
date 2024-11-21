package router

import (
	"nbientry/web/common"
	"net/http"
)

func HomeRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if user.Email == "" {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}
	// fail.Render(w, r, pages.Home(user))
}
