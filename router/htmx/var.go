package htmx

import (
	"fmt"
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/common"
	"nbientry/web/variable"
	"net/http"
	"strconv"
)

func VariableRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, r, user, fail.ADMIN) {
		return
	}

	if len(path) != 1 {
		http.Error(w, "Bad request\nformat: /htmx/var/:VariableID", http.StatusBadRequest)
		return
	}

	v_id, err := strconv.Atoi(path[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v := web.GetFirst(variable.Variable{Id: v_id})
	if v.Id == 0 {
		http.Error(w, "Variable not found", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		Variable_PUT(w, r, user, v)

	case "PATCH":
		fail.Render(w, r, v.RenderInEditor())

	case "GET":
		fail.Render(w, r, v.RenderInViewer())

	case "POST":
		Variable_POST(w, r, user, v)

	case "DELETE":
		parent_id := v.ParentId
		if parent_id == 0 {
			http.Error(w, "Cannot delete root variable", http.StatusBadRequest)
			return
		}

		parent := web.GetFirst(variable.Variable{Id: parent_id})
		if parent.Id == 0 {
			http.Error(w, "Parent variable not found", http.StatusBadRequest)
			return
		}

		v.Delete()
		w.Header().Set("HX-Retarget", fmt.Sprintf("#variable-%d", parent.Id))
		fail.Render(w, r, parent.RenderInViewer())

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Variable_PUT(w http.ResponseWriter, r *http.Request, user common.User, v variable.Variable) {
	if fail.Form(w, r) {
		return
	}

	t := r.Form.Get("type")

	new_v := variable.Variable{
		ParentId: v.Id,
		Name:     "New [" + t + "]",
		Type:     t,
	}

	_, err := new_v.New().ToTypedEntry()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	new_v.Id = 0
	web.Save(&new_v)
	fail.Render(w, r, v.RenderInViewer())
}

func Variable_POST(w http.ResponseWriter, r *http.Request, user common.User, v variable.Variable) {
	if fail.Form(w, r) {
		return
	}

	v.Name = r.Form.Get("name")
	if t := r.Form.Get("type"); t != "" {
		v.Type = t
	}
	v.Description = r.Form.Get("description")
	v.Suffix = r.Form.Get("suffix")

	_, err := v.New().ToTypedEntry()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	web.Save(&v)
	fail.Render(w, r, v.RenderInViewer())
}
