package router

import (
	"nbientry/web/common"
	"net/http"
	"os"
	"strings"
)

func Router(w http.ResponseWriter, r *http.Request) {
	user := common.CookieAuth(w, r)

	path := strings.Split(r.URL.Path[len("/"):], "/")
	if len(path) < 2 {
		path = append(path, "")
	}

	switch path[0] {
	case "src":
		HandleSource(w, r)

	case "htmx":
		HtmxRouter(w, r, user, path[1:])

	case "":
		HomeRouter(w, r, user, path[1:])

	case "user":
		UserRouter(w, r, user, path[1:])

	case "admin":
		AdminRouter(w, r, user, path[1:])

	default:
		w.Header().Set("X-Redirect-Reason", "404: /"+path[0])
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

func HandleSource(w http.ResponseWriter, r *http.Request) {
	file := "." + r.URL.Path
	stat, err := os.Stat(file)
	if err != nil || stat.Size() == 0 {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	// Cache Images
	if strings.HasPrefix(file, "./src/img/upload/") {
		w.Header().Set("Cache-Control", "max-age=31536000")
	}

	http.ServeFile(w, r, file)
}
