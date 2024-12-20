package htmx

import (
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/common"
	"nbientry/web/notif"
	"net/http"
	"strconv"
)

func StageRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, r, user, fail.ADMIN) {
		return
	}

	stage_id, err := strconv.Atoi(path[0])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	stage := web.GetFirst(notif.Stage{Id: stage_id})

	switch r.Method {
	case "PUT":
		Stage_PUT(w, r, user, path, stage)

	case "PATCH":
		fail.Render(w, r, stage.RenderHead_Edit())

	case "GET":
		fail.Render(w, r, stage.RenderHead_View())

	case "DELETE":
		stage.Delete()
		w.Header().Set("HX-Refresh", "true")

	case "POST":
		if fail.Form(w, r) {
			return
		}
		was_empty := stage.Name == ""
		stage.Name = r.Form.Get("name")
		stage.Description = r.Form.Get("description")
		web.Save(&stage)
		if was_empty {
			w.Header().Set("HX-Refresh", "true")
		}
		fail.Render(w, r, stage.RenderHead_View())

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func Stage_PUT(w http.ResponseWriter, r *http.Request, user common.User, path []string, stage notif.Stage) {
	if len(path) != 2 {
		http.Error(w, "Bad request\nformat: /htmx/stage/:StageID/:Direction (direction can be inc, dec, new)", http.StatusBadRequest)
		return
	}
	w.Header().Set("HX-Retarget", "#stage-table")

	switch path[1] {
	case "inc":
		stage.Increment()
		fail.Render(w, r, notif.RenderStageTable())

	case "dec":
		stage.Decrement()
		fail.Render(w, r, notif.RenderStageTable())

	case "new":
		_ = stage.New()
		fail.Render(w, r, notif.RenderStageTable())
	}
}
