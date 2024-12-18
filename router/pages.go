package router

import (
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/account"
	"nbientry/web/common"
	"nbientry/web/notif"
	"nbientry/web/pages"
	"net/http"
	"strconv"
)

func HomeRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if user.Email == "" {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}
	fail.Render(w, r, pages.Home(user))
}

func AdminRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, r, user, fail.ADMIN) {
		return
	}

	switch path[0] {
	case "user":
		fail.Render(w, r, account.UserList(user))

	case "stage":
		StageRouter(w, r, user, path[1:])

	default:
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func StageRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, r, user, fail.ADMIN) {
		return
	}

	if len(path) == 0 {
		fail.Render(w, r, notif.RenderStageTablePage(user))
		return
	}

	if len(path) != 1 {
		w.Header().Set("X-Redirect-Reason", "Bad request: /admin/stage/:ID")
		http.Redirect(w, r, "/admin/stage", http.StatusTemporaryRedirect)
		return
	}

	id, err := strconv.Atoi(path[0])
	if err != nil {
		w.Header().Set("X-Redirect-Reason", "Not a number; /admin/stage/:ID")
		http.Redirect(w, r, "/admin/stage", http.StatusTemporaryRedirect)
		return
	}

	stage := web.GetFirst(notif.Stage{Id: id})
	if stage.Id != id {
		w.Header().Set("X-Redirect-Reason", "Not found; /admin/stage/:ID")
		http.Redirect(w, r, "/admin/stage", http.StatusTemporaryRedirect)
		return
	}

	fail.Render(w, r, stage.RenderPage(user))
}
