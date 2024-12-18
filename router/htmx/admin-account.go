package htmx

import (
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/common"
	"net/http"
	"strconv"
)

func AdminAccountRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, r, user, fail.ADMIN) {
		return
	}

	if len(path) != 1 {
		http.Error(w, "Bad request\nformat: /htmx/admin-account/:UserID", http.StatusBadRequest)
		return
	}

	u_id, err := strconv.Atoi(path[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u := web.GetFirst(common.User{Id: u_id})
	if u.Id == 0 {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		fail.Render(w, r, u.RenderChip_PostView(user))

	case "PATCH":
		u.Admin = !u.Admin
		web.Save(&u)
		fail.Render(w, r, u.RenderChip_PostView(user))

	case "DELETE":
		u.Password = u.HashPassword("password")
		web.Save(&u)
		fail.Render(w, r, u.RenderChip_PostView(user))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
