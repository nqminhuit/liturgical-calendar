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

func WeekPass(filename string) {
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
