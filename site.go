package main

import (
	"fmt"
	"nbientry/router"
	"net/http"
)

func main() {
	http.HandleFunc("/", router.Router)

	fmt.Println("Listening on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
