package calendar

import (
	"context"
	"fmt"
	"time"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) GetCalendar(ctx context.Context, month int, year int) (Calendar, error) {
	// using current day, calculate the weekday of first day
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	firstWeekday := firstDay.Weekday()

	// We know the 28th is the earliest the last day of the month will ever be (February)
	earliestLastDay := time.Date(year, time.Month(month), 28, 0, 0, 0, 0, time.Now().Location())
	lastDay := s.findEndOfMonth(earliestLastDay)

	weeksInMonth := (lastDay.Day() + int(firstWeekday)) / 7

	hasRemainderDays := (lastDay.Day()+int(firstWeekday))%7 > 0

	if hasRemainderDays {
		weeksInMonth += 1
	}

	if weeksInMonth < 4 || weeksInMonth > 6 {
		return Calendar{}, fmt.Errorf("Calendar can only have between 4 and 6 weeks but has %d", weeksInMonth)
	}

	calendarDates := make([][]uint8, weeksInMonth)

	dayPointer := firstDay.Day() + int(firstWeekday)
	dayToAdd := 1
	for idxWeek := range calendarDates {
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
		calendarDates[idxWeek] = week
	}

	if dayToAdd <= lastDay.Day() {
		return Calendar{}, fmt.Errorf("not all days in month added to calendar")
	}

	header := fmt.Sprintf("%s %d", firstDay.Month().String(), firstDay.Year())

	calendar := Calendar{
		Header: header,
		Dates:  calendarDates,
	}

	if err := calendar.Validate(); err != nil {
		return Calendar{}, fmt.Errorf("calendar failed validation: %w", err)
	}

	return calendar, nil
}

func (s Service) findEndOfMonth(date time.Time) time.Time {
	nextDay := date.AddDate(0, 0, 1)
	if nextDay.Month() > date.Month() {
		return date
	}
	return s.findEndOfMonth(nextDay)
}

type Calendar struct {
	Header string
	Dates  [][]uint8
}

// Validate asserts that a calendar is enclosed with empty dates, that dates are linear, and that weeks are non empty
func (c Calendar) Validate() error {
	dayPointer := 1
	beforeFirstValidDay := true
	afterValidLastDay := false
	for _, week := range c.Dates {
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

func (c Calendar) String() [][]string {
	output := make([][]string, len(c.Dates))

	for iWeek, week := range c.Dates {
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
