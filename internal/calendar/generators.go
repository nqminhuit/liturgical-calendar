package calendar

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

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

// ordinaryWeekPass traverses through calendar map,
// filter only season == "ordinary" to set week_of_season.
// Ordinary in a liturgical calendar has 2 segments:
//
// 1. after Christmas, until Ash Wednesday (the first wednesday of Lent with week_of_season = 0)
//
// 2. after Pentecost
//
// week_of_season from the 1st segment will increase by 1 on every Sunday.
// week_of_season from the 2nd segment will have to traverse backward:
// from the first Sunday of Advent Season back to 1 Sunday is the 34th week of Ordinary Season,
// then keep going back and set the week_of_season.
//
// There are always 34 weeks of Ordinary season, never more, never less.
func ordinaryWeekPass(calendar map[string]DayInfo) {
	// ---- sort dates ----
	dates := make([]string, 0, len(calendar))
	for d := range calendar {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	// =====================================
	// FIRST SEGMENT (forward)
	// =====================================
	week := 0
	started := false
	for _, d := range dates {
		day := calendar[d]
		if day.Season == SeasonOrdinary {
			if !started {
				started = true
				week = 1
			} else if day.Weekday == "sun" {
				week++
			}

			day.WeekOfSeason = week
			calendar[d] = day
		} else if started {
			// stop when leaving first ordinary segment
			break
		}
	}

	// =====================================
	// FIND ADVENT START
	// =====================================
	adventStart := -1
	for i, d := range dates {
		if calendar[d].Season == SeasonAdvent {
			adventStart = i
			break
		}
	}
	if adventStart == -1 {
		return
	}

	// =====================================
	// SECOND SEGMENT (backward)
	// =====================================
	week = 34

	for i := adventStart - 1; i >= 0; i-- {
		d := dates[i]
		day := calendar[d]

		if day.Season != SeasonOrdinary {
			break
		}

		day.WeekOfSeason = week
		calendar[d] = day

		if day.Weekday == "sun" {
			week--
			if week == 0 {
				break
			}
		}
	}
}

func specialSeasonPass(calendar map[string]DayInfo) {
	// sort dates
	dates := make([]string, 0, len(calendar))
	for k := range calendar {
		dates = append(dates, k)
	}
	sort.Strings(dates)

	// counters
	var (
		adventWeek    int
		christmasWeek int
		lentWeek      int
		easterWeek    int
		prevSeason    Season = ""
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
			// ignore here
		}

		calendar[key] = day
		prevSeason = day.Season
	}

}

func lectionaryPass(calendar map[string]DayInfo) {
	for d, day := range calendar {
		day.Lectionary = day.lectionaryKey()
		calendar[d] = day
	}
}

func WeekPass(filename string) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	calendar := map[string]DayInfo{}
	if err := json.Unmarshal(raw, &calendar); err != nil {
		panic(err)
	}

	specialSeasonPass(calendar)
	ordinaryWeekPass(calendar)
	lectionaryPass(calendar)

	out, err := json.MarshalIndent(calendar, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(filename, out, 0644)
	fmt.Println("Generated:", filename)
}
