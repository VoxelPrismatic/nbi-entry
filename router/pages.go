package router

import (
	"nbientry/router/fail"
	"nbientry/web/common"
	notif "nbientry/web/notification"
	"nbientry/web/pages"
	"net/http"
)

func HomeRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if user.Email == "" {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}
	fail.Render(w, r, pages.Home(user))
}

func AdminRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, user, fail.ADMIN) {
		return
	}

	switch path[0] {
	case "notif":
		fail.Render(w, r, notif.NotifTable(user))
	}
}
