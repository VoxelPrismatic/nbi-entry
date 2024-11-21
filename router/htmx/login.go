package htmx

import (
	"fmt"
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/account"
	"nbientry/web/common"
	"net/http"
)

func LoginRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	switch r.Method {
	case "POST":
		Login_POST(w, r, user)

	case "PATCH":
		Login_PATCH(w, r, user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ValidateUser(r *http.Request, user common.User, create bool) (common.User, error) {
	err := r.ParseForm()
	if err != nil {
		return user, err
	}

	user = common.User{
		Email: r.Form.Get("username"),
	}

	password := r.Form.Get("password")
	err = user.ValidatePassword(password)
	if err != nil {
		return user, err
	}

	if !user.Exists() {
		check_pw := r.Form.Get("check-pw")
		if password != check_pw {
			return user, fmt.Errorf("passwords don't match")
		}

		user.Password = user.HashPassword(password)

		if create {
			var count int64
			web.Db().Model(&common.User{}).Count(&count)
			if count == 0 {
				user.Admin = true
			}

			fmt.Println(user)

			web.Save(&user)
		} else {
			return user, fmt.Errorf("-")
		}
	}

	return user, nil
}

func Login_POST(w http.ResponseWriter, r *http.Request, user common.User) {
	msg := ""
	user, err := ValidateUser(r, user, false)
	if err != nil {
		msg = err.Error()
	}

	if msg == "-" {
		w.Header().Set("HX-Retarget", "#login-btn")
		w.Header().Set("HX-Reswap", "outerHTML")
		fail.Render(w, r, account.LoginBtn("", !user.Exists()))
	} else if msg == "passwords don't match" {
		w.Header().Set("HX-Retarget", "#login-btn")
		w.Header().Set("HX-Reswap", "outerHTML")
		fail.Render(w, r, account.LoginBtn(msg, !user.Exists()))
	} else {
		fail.Render(w, r, account.LoginBox(msg, !user.Exists()))
	}
}

func Login_PATCH(w http.ResponseWriter, r *http.Request, user common.User) {
	user, err := ValidateUser(r, user, true)
	if err != nil {
		fail.Render(w, r, account.LoginBox(err.Error(), !user.Exists()))
	}

	err = user.GenerateJWT()
	if err != nil {
		fail.Render(w, r, account.Login(err.Error(), !user.Exists()))
	}

	http.SetCookie(w, user.Cookie())
	w.Header().Set("HX-Redirect", "/user/")
}
