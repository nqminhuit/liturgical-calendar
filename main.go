package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Season string

var loc *time.Location

func init() {
	var err error
	if loc, err = time.LoadLocation("Asia/Ho_Chi_Minh"); err != nil {
		panic(err)
	}
}

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

// ---------------- Easter ----------------

func Easter(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func startOfAdvent(year int) time.Time {
	// Sunday between Nov 27 and Dec 3
	start := time.Date(year, 11, 27, 0, 0, 0, 0, loc)
	offset := (7 - int(start.Weekday())) % 7
	return start.AddDate(0, 0, offset)
}

func epiphany(year int) time.Time {
	d := time.Date(year, 1, 2, 0, 0, 0, 0, loc) // Jan 2
	offset := (7 - int(d.Weekday())) % 7        // move forward to Sunday
	return d.AddDate(0, 0, offset)
}

func baptismOfLord(year int) time.Time {
	e := epiphany(year)
	nextSunday := e.AddDate(0, 0, (7-int(e.Weekday()))%7)
	if nextSunday.Equal(e) {
		nextSunday = e.AddDate(0, 0, 7)
	}
	if e.Day() >= 7 { // edge case: if Epiphany Jan 7 or 8 → Monday
		return e.AddDate(0, 0, 1)
	}
	return nextSunday
}

// Liturgical year starts at Advent
func liturgicalYear(d time.Time) int {
	advent := startOfAdvent(d.Year())
	if !d.Before(advent) {
		return d.Year() + 1
	}
	return d.Year()
}

// Sunday cycle repeats every 3 years, starting with A in 2008 (first Sunday of Advent)
func sundayCycle(d time.Time) string {
	ly := liturgicalYear(d)
	switch (ly - 2008) % 3 {
	case 0:
		return "A"
	case 1:
		return "B"
	default:
		return "C"
	}
}

func weekdayCycle(year int) string {
	if year%2 == 0 {
		return "II"
	}
	return "I"
}

func seasonOf(d time.Time, ash, easter, pentecost, advent time.Time) Season {
	year := d.Year()
	christmasCurrent := time.Date(year, 12, 25, 0, 0, 0, 0, loc)
	christmasPrev := time.Date(year-1, 12, 25, 0, 0, 0, 0, loc)
	baptism := baptismOfLord(year)

	switch {
	// Christmas season from previous Dec 25 -> Baptism
	case !d.Before(christmasPrev) && d.Before(baptism.AddDate(0, 0, 1)):
		return SeasonChristmas

	// Christmas season from this Dec 25 -> end of year
	case !d.Before(christmasCurrent):
		return SeasonChristmas

	case !d.Before(ash) && d.Before(easter):
		return SeasonLent

	case !d.Before(easter) && !d.After(pentecost):
		return SeasonEaster

	case !d.Before(advent) && d.Before(christmasCurrent):
		return SeasonAdvent

	default:
		return SeasonOrdinary
	}
}

func GenerateYear(year int) map[string]DayInfo {
	result := map[string]DayInfo{}

	easter := Easter(year)
	ash := easter.AddDate(0, 0, -46)
	pentecost := easter.AddDate(0, 0, 49)
	advent := startOfAdvent(year)

	start := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	end := time.Date(year, 12, 31, 0, 0, 0, 0, loc)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		season := seasonOf(d, ash, easter, pentecost, advent)
		key := d.Format("2006-01-02")
		weekday := strings.ToLower(d.Weekday().String()[:3])

		result[key] = DayInfo{
			Season:       season,
			SundayCycle:  sundayCycle(d),
			WeekdayCycle: weekdayCycle(year),
			Weekday:      weekday,
		}
	}

	return result
}

func weekPass(filename string) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	calendar := map[string]DayInfo{}
	if err := json.Unmarshal(raw, &calendar); err != nil {
		panic(err)
	}

	// sort dates
	dates := make([]string, 0, len(calendar))
	for k := range calendar {
		dates = append(dates, k)
	}
	sort.Strings(dates)

	// counters
	var (
		adventWeek             int
		christmasWeek          int
		lentWeek               int
		easterWeek             int
		ordinaryWeek           int
		prevSeason             Season = ""
		lastOrdinaryBeforeLent int
	)

	for _, key := range dates {
		day := calendar[key]

		// detect sunday
		d, _ := time.Parse("2006-01-02", key)
		isSunday := d.Weekday() == time.Sunday

		switch day.Season {

		case SeasonAdvent:
			if prevSeason != SeasonAdvent {
				adventWeek = 1
			} else if isSunday {
				adventWeek++
			}
			day.WeekOfSeason = adventWeek

		case SeasonChristmas:
			if prevSeason == SeasonAdvent {
				christmasWeek = 0
			} else if prevSeason != SeasonChristmas {
				christmasWeek = 1
			} else if isSunday {
				christmasWeek++
			}
			day.WeekOfSeason = christmasWeek

		case SeasonLent:
			if prevSeason == SeasonOrdinary {
				lastOrdinaryBeforeLent = ordinaryWeek
			}
			if prevSeason != SeasonLent {
				lentWeek = 0
			} else if isSunday {
				lentWeek++
			}
			day.WeekOfSeason = lentWeek

		case SeasonEaster:
			if prevSeason != SeasonEaster {
				easterWeek = 1
			} else if isSunday {
				easterWeek++
			}
			day.WeekOfSeason = easterWeek

		case SeasonOrdinary:
			// IMPORTANT:
			// Ordinary continues after Easter → do NOT reset
			if prevSeason != SeasonOrdinary {
				if ordinaryWeek == 0 {
					ordinaryWeek = 1
				}
				if prevSeason == SeasonEaster {
					ordinaryWeek = lastOrdinaryBeforeLent + 2
				}
			} else if isSunday {
				ordinaryWeek++
			}
			day.WeekOfSeason = ordinaryWeek
		}

		calendar[key] = day
		prevSeason = day.Season
	}

	out, err := json.MarshalIndent(calendar, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(filename, out, 0644)
	fmt.Println("Generated:", filename)
}

// Usage: go run main.go <year>
func main() {
	if len(os.Args) < 2 {
		panic("Usage: go run main.go <year>")
	}
	yearStr := os.Args[1]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		panic("Invalid year: " + yearStr)
	}
	cal := GenerateYear(year)

	b, err := json.MarshalIndent(cal, "", "  ")
	if err != nil {
		panic(err)
	}
	filename := fmt.Sprintf("liturgical-calendar-%d.json", year)
	os.WriteFile(filename, b, 0644)
	weekPass(filename)
}
