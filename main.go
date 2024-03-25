package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/alexander-littleton/go-htmx-project/pkg/pages"
)

func main() {
	now := time.Now()
	// using current day, calculate the weekday of first day
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstWeekday := firstDay.Weekday()

	// We know the 28th is the earliest the last day of the month will ever be (February)
	earliestLastDay := time.Date(now.Year(), now.Month(), 28, 0, 0, 0, 0, now.Location())
	lastDay := findEndOfMonth(earliestLastDay)

	weeksInMonth := (lastDay.Day() + int(firstWeekday)) / 7

	hasRemainderDays := (lastDay.Day()+int(firstWeekday))%7 > 0

	if hasRemainderDays {
		weeksInMonth += 1
	}

	if weeksInMonth < 4 || weeksInMonth > 6 {
		fmt.Printf("ERR: Calendar can only have between 4 and 6 weeks but has %v", weeksInMonth)
		return
	}

	calendar := make(calendar, weeksInMonth)

	dayPointer := firstDay.Day() + int(firstWeekday)
	dayToAdd := 1
	for idxWeek := range calendar {
		week := make([]uint8, 7)

		for idxDay := range week {
			numDay := idxDay + 1 + idxWeek*7
			if numDay == dayPointer {
				week[idxDay] = uint8(dayToAdd)
				dayToAdd++
				dayPointer++
			}

			if dayToAdd > lastDay.Day() {
				break
			}
		}
		calendar[idxWeek] = week
	}

	if dayToAdd <= lastDay.Day() {
		fmt.Println("ERR: Not all days in month added to calendar")
		return
	}

	if err := calendar.Validate(); err != nil {
		fmt.Println("ERR: Calendar failed validation: ", err.Error())
	}

	header := fmt.Sprintf("%s %d", now.Month().String(), now.Year())

	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(pages.Home(header, calendar.String())))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", mux)
}

func findEndOfMonth(date time.Time) time.Time {
	nextDay := date.AddDate(0, 0, 1)
	if nextDay.Month() > date.Month() {
		return date
	}
	return findEndOfMonth(nextDay)
}

type calendar [][]uint8

// Validate asserts that a calendar is enclosed with empty dates, that dates are linear, and that weeks are non empty
func (c calendar) Validate() error {
	dayPointer := 1
	beforeFirstValidDay := true
	afterValidLastDay := false
	for _, week := range c {
		for i, day := range week {
			if afterValidLastDay && day == 0 {
				continue
			}
			if afterValidLastDay && (day < 28 || day > 31) && day != 0 {
				return fmt.Errorf("invalid date %d after last date of month", day)
			}
			if day > 0 {
				beforeFirstValidDay = false
			}
			if !beforeFirstValidDay && int(day) != dayPointer {
				return fmt.Errorf("expected calendar day %d but received %d", dayPointer, day)
			}
			if day > 0 {
				dayPointer++
			}
			if day == 28 {
				afterValidLastDay = true
			}
			if i == 6 && beforeFirstValidDay {
				return fmt.Errorf("first week is empty")
			}
			if i == 0 && day == 0 && afterValidLastDay {
				return fmt.Errorf("extra week detected")
			}
		}
	}
	return nil
}

func (c calendar) String() [][]string {
	output := make([][]string, len(c))

	for iWeek, week := range c {
		outputWeek := make([]string, 7)
		for iDay, day := range week {
			if day != 0 {
				outputWeek[iDay] = fmt.Sprint(day)
			}
		}
		output[iWeek] = outputWeek
	}

	return output
}
