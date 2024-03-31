package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexander-littleton/go-htmx-project/pkg/domain/calendar"
)

func main() {

	mux := http.NewServeMux()

	calendarSvc := calendar.NewService()
	calendarCtrl := calendar.NewController(calendarSvc)
	calendar.InitRoutes(mux, calendarCtrl)

	mux.Handle("GET /", templ.Handler(Home()))
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", mux)
}
