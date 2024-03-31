package calendar

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexander-littleton/go-htmx-project/pkg/domain/calendar/pages"
)

type CalendarController interface {
	getCalendar(w http.ResponseWriter, req *http.Request)
}

func InitRoutes(mux *http.ServeMux, controller CalendarController) {
	mux.HandleFunc("GET /calendar", controller.getCalendar)
}

type CalendarService interface {
	GetCalendar(ctx context.Context, month, year int) (Calendar, error)
}

type Controller struct {
	calendarService CalendarService
}

func NewController(calendarService CalendarService) Controller {
	return Controller{
		calendarService: calendarService,
	}
}

func (c Controller) getCalendar(w http.ResponseWriter, req *http.Request) {
	monthYear := req.URL.Query().Get("monthYear")
	split := strings.Split(monthYear, "-")
	month, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Println("failed to parse month from request")
		return
	}
	year, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("failed to parse year from request")
		return
	}

	ctx := context.Background()
	calendar, err := c.calendarService.GetCalendar(ctx, month, year)
	if err != nil {
		fmt.Println(err.Error())
	}

	pages.Calendar(calendar.Header, calendar.String()).Render(ctx, w)
}
