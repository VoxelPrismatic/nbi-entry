package htmx

import (
	"fmt"
	"nbientry/router/fail"
	"nbientry/web"
	"nbientry/web/common"
	"nbientry/web/notif"
	"net/http"
	"strconv"
)

func NotifRouter(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if fail.Auth(w, r, user, fail.ADMIN) {
		return
	}

	switch r.Method {
	case "PUT":
		Notif_PUT(w, r, user, path)

	case "PATCH":
		Notif_PATCH(w, r, user, path)

	case "POST":
		Notif_POST(w, r, user, path)

	case "GET":
		Notif_GET(w, r, user, path)

	case "DELETE":
		Notif_DELETE(w, r, user, path)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func notif_parse_url(path []string, has_stage bool, has_app_seg bool) (notif.Stage, notif.ApplicationSegment, error) {
	flags := 0
	app_seg_idx := -1
	stage_idx := -1
	if has_stage {
		stage_idx++
		flags |= 1
	}
	if has_app_seg {
		app_seg_idx++
		stage_idx++
		flags |= 2
	}

	stage := notif.Stage{}
	app_seg := notif.ApplicationSegment{}

	switch flags {
	case 0:
		if len(path) != 0 {
			return stage, app_seg, fmt.Errorf("Bad request\nformat: /htmx/notif/")
		}
	case 1:
		if len(path) != 1 {
			return stage, app_seg, fmt.Errorf("Bad request\nformat: /htmx/notif/:StageID")
		}
	case 2:
		if len(path) != 1 {
			return stage, app_seg, fmt.Errorf("Bad request\nformat: /htmx/notif/:AppSegID")
		}
	case 3:
		if len(path) != 2 {
			return stage, app_seg, fmt.Errorf("Bad request\nformat: /htmx/notif/:AppSegID/:StageID")
		}
	}

	if has_stage {
		stage_id, err := strconv.Atoi(path[stage_idx])
		if err != nil {
			return stage, app_seg, fmt.Errorf("Invalid stage ID")
		}
		stage = web.GetFirst(notif.Stage{Id: stage_id})
		if stage.Id != stage_id {
			return stage, app_seg, fmt.Errorf("Stage not found")
		}
	}

	if has_app_seg {
		app_seg_id, err := strconv.Atoi(path[app_seg_idx])
		if err != nil {
			return stage, app_seg, fmt.Errorf("Invalid application segment ID")
		}
		app_seg = web.GetFirst(notif.ApplicationSegment{Id: app_seg_id})
		if app_seg.Id != app_seg_id {
			return stage, app_seg, fmt.Errorf("Application segment not found")
		}
	}

	return stage, app_seg, nil
}

func Notif_PUT(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	if len(path) == 2 {
		stage, app_seg, err := notif_parse_url(path, true, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if fail.Form(w, r) {
			return
		}

		u_id, err := strconv.Atoi(r.Form.Get("user"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u := web.GetFirst(common.User{Id: u_id})
		if u.Id == 0 {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}

		n := notif.Notification{AppSegId: app_seg.Id, StageId: stage.Id, UserId: u.Id}
		web.Save(&n)

		fail.Render(w, r, stage.RenderApplicationSegment_View(app_seg))
		return
	}

	stage, _, err := notif_parse_url(path, true, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := notif.ApplicationSegment{Name: "New Application Segment", Color: "rose"}
	web.Save(&app)

	fail.Render(w, r, stage.RenderApplicationSegment_Template(app))
}

func Notif_PATCH(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	stage, app_seg, err := notif_parse_url(path, true, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fail.Render(w, r, stage.RenderApplicationSegment_Edit(app_seg))
}

func Notif_POST(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	stage, app_seg, err := notif_parse_url(path, true, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if fail.Form(w, r) {
		return
	}

	app_seg.Name = r.Form.Get("name")
	app_seg.Color = r.Form.Get("color")
	web.Save(&app_seg)

	fail.Render(w, r, stage.RenderApplicationSegment_View(app_seg))
}

func Notif_GET(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	stage, app_seg, err := notif_parse_url(path, true, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fail.Render(w, r, stage.RenderApplicationSegment_View(app_seg))
}

func Notif_DELETE(w http.ResponseWriter, r *http.Request, user common.User, path []string) {
	fmt.Println(path)
	fmt.Println(len(path))
	if len(path) == 3 {
		u_id, err := strconv.Atoi(path[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u := web.GetFirst(common.User{Id: u_id})
		if u.Id == 0 {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}

		stage, app_seg, err := notif_parse_url(path[:2], true, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		n := notif.Notification{AppSegId: app_seg.Id, StageId: stage.Id, UserId: u.Id}
		web.Db().Model(&notif.Notification{}).Where(&n).Delete(&n)

		w.Header().Set("HX-Retarget", fmt.Sprintf("#app-%d", app_seg.Id))
		fail.Render(w, r, stage.RenderApplicationSegment_View(app_seg))
		return
	}

	_, app_seg, err := notif_parse_url(path, false, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, n := range web.GetSorted(notif.Notification{AppSegId: app_seg.Id}, "stage_id ASC") {
		web.Db().Model(&notif.Notification{}).Where(&n).Delete(&n)
	}
	web.Db().Delete(&app_seg)
}
