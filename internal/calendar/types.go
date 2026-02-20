package calendar

import "time"

type Season string

var loc *time.Location

const (
	SeasonAdvent    Season = "advent"
	SeasonChristmas Season = "christmas"
	SeasonLent      Season = "lent"
	SeasonEaster    Season = "easter"
	SeasonOrdinary  Season = "ordinary"
)

type DayInfo struct {
	Season       Season `json:"season"`
	SundayCycle  string `json:"sunday_cycle"`  // A B C
	WeekdayCycle string `json:"weekday_cycle"` // I II
	Weekday      string `json:"weekday"`
	WeekOfSeason int    `json:"week_of_season"`
}

func init() {
	var err error
	if loc, err = time.LoadLocation("Asia/Ho_Chi_Minh"); err != nil {
		panic(err)
	}
}
