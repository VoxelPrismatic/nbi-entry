package router

import (
	"fmt"
	"nbientry/router/fail"
	"nbientry/router/htmx"
	"nbientry/web"
	"nbientry/web/account"
	"nbientry/web/common"
	"net/http"
	"strconv"
)

func UserRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if len(path) == 0 {
		path = append(path, "")
	}

	switch path[0] {
	case "login":
		LoginRouter(w, r, user, path[1:])

	case "view":
		UserViewRouter(w, r, user, path[1:])

	default:
		if user.Email == "" {
			http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/user/view/%d", user.Id), http.StatusTemporaryRedirect)
	}

}

func LoginRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	switch r.Method {
	case "GET":
		w.Header().Set("Set-Cookie", "jwt=nil; path=/")
		fail.Render(w, r, account.Login("email too short", false))

	case "POST":
		Login_POST(w, r, user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Login_POST(w http.ResponseWriter, r *http.Request, user common.User) {
	if fail.Form(w, r) {
		return
	}

	user, err := htmx.ValidateUser(r, user, true)
	if err != nil {
		fail.Render(w, r, account.Login(err.Error(), !user.Exists()))
	}

	err = user.GenerateJWT()
	if err != nil {
		fail.Render(w, r, account.Login(err.Error(), !user.Exists()))
	}

	http.SetCookie(w, user.Cookie())
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func UserViewRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	switch r.Method {
	case "GET":
		UserView_GET(w, r, user, path)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UserView_GET(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if len(path) == 0 {
		http.Redirect(w, r, "/user/", http.StatusTemporaryRedirect)
		return
	}

	if len(path) != 1 {
		http.Error(w, "Bad request\nformat: /user/view/ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path[0])
	if err != nil {
		http.Redirect(w, r, "/user/", http.StatusTemporaryRedirect)
		return
	}

	target := web.GetFirst(common.User{Id: id})
	if target.Email == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	fail.Render(w, r, account.UserPage(user, target))
}
