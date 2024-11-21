package htmx

import (
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/common"
	"net/http"
	"strconv"
)

func AccountRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if len(path) != 1 {
		http.Error(w, "Bad request\nformat: /htmx/account/ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "POST":
		Account_POST(w, r, user, path)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func Account_POST(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	id, err := strconv.Atoi(path[0])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if fail.Form(w, r) {
		return
	}

	target := web.GetFirst(common.User{Id: id})
	if target.Email == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if target.Id != user.Id {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user.FullName = r.Form.Get("name")
	user.Email = r.Form.Get("email")
	user.Phone = r.Form.Get("phone")
	user.Image = r.Form.Get("image")

	if !common.ValidNames.MatchString(user.Email) {
		w.Header().Set("HX-Retarget", "#email")
		fail.Render(w, r, user.RenderEmail_Edit("Invalid email"))
		return
	}

	pass := r.Form.Get("password")
	if pass != "" {
		err = user.ValidatePassword(pass)
		if err != nil {
			if err.Error() != "bad credentials" {
				w.Header().Set("HX-Retarget", "#password")
				fail.Render(w, r, user.RenderPassword_Edit(err.Error()))
				return
			}

			user.Password = user.HashPassword(pass)
		}
	}

	web.Save(user)
	fail.Render(w, r, user.RenderUser_View())
}
