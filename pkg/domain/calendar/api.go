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
	getBookingForm(w http.ResponseWriter, req *http.Request)
}

func InitRoutes(mux *http.ServeMux, controller CalendarController) {
	mux.HandleFunc("GET /calendar", controller.getCalendar)
	mux.HandleFunc("GET /bookingForm", controller.getBookingForm)
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

	selectedDay := req.URL.Query().Get("selectedDay")

	pages.Calendar(month, year, calendar.String(), selectedDay).Render(ctx, w)
}

func (c Controller) getBookingForm(w http.ResponseWriter, req *http.Request) {

}
