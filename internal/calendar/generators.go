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

		dayInfo := DayInfo{
			Season:      season,
			SundayCycle: sundayCycle(d),
			Weekday:     weekday,
		}
		// only set weekday cycle for ordinary season
		if season == SeasonOrdinary {
			dayInfo.WeekdayCycle = weekdayCycle(year)
		}
		result[key] = dayInfo
	}

	return result
}

func applyWeekNumbers(calendar map[string]DayInfo) {
	dates := make([]string, 0, len(calendar))
	for d := range calendar {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	var (
		adventWeek      int
		christmasWeek   int
		lentWeek        int
		easterWeek      int
		ordWeek         int
		ordStarted      bool
		prevSeason      Season = ""
		pastLent        bool
		adventStart     int = -1
	)

	for i, key := range dates {
		day := calendar[key]
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
			if adventStart == -1 {
				adventStart = i
			}

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
			pastLent = true

		case SeasonEaster:
			if prevSeason != SeasonEaster {
				easterWeek = 1
			} else if isSunday {
				easterWeek++
			}
			day.WeekOfSeason = easterWeek

		case SeasonOrdinary:
			if !pastLent {
				if !ordStarted {
					ordStarted = true
					ordWeek = 1
				} else if isSunday {
					ordWeek++
				}
				day.WeekOfSeason = ordWeek
			}
		}

		calendar[key] = day
		prevSeason = day.Season
	}

	if adventStart == -1 {
		return
	}

	ordWeek = 34
	for i := adventStart - 1; i >= 0; i-- {
		key := dates[i]
		day := calendar[key]

		if day.Season != SeasonOrdinary {
			break
		}

		day.WeekOfSeason = ordWeek
		calendar[key] = day

		if day.Weekday == "sun" {
			ordWeek--
			if ordWeek == 0 {
				break
			}
		}
	}
}

func lectionaryPass(calendar map[string]DayInfo) {
	for d, day := range calendar {
		day.Lectionary = day.lectionaryKey()
		calendar[d] = day
	}
}

func loadLectionary(resourceDir string) map[string]map[string]any {
	raw, err := os.ReadFile(resourceDir + "/lectionary.json")
	if err != nil {
		panic(err)
	}
	var lec struct {
		Readings map[string]map[string]any `json:"readings"`
	}
	if err := json.Unmarshal(raw, &lec); err != nil {
		panic(err)
	}
	return lec.Readings
}

func applyNames(calendar map[string]DayInfo, readings map[string]map[string]any) {
	for date, day := range calendar {
		if day.Name != "" {
			continue
		}
		if day.Lectionary == "" {
			continue
		}
		if rd, ok := readings[day.Lectionary]; ok && rd != nil {
			if nI, ok := rd["name"]; ok {
				if nStr, ok2 := nI.(string); ok2 && nStr != "" {
					day.Name = nStr
					calendar[date] = day
				}
			}
		}
	}
}

func GenerateCalendar(year int, resourceDir string) map[string]DayInfo {
	cal := GenerateYear(year)
	applyWeekNumbers(cal)
	lectionaryPass(cal)
	readings := loadLectionary(resourceDir)
	applyNames(cal, readings)
	return cal
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

	applyWeekNumbers(calendar)
	lectionaryPass(calendar)

	out, err := json.MarshalIndent(calendar, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filename, out, 0644); err != nil {
		panic(err)
	}
	fmt.Println(filename)
}

// another pass to set the name of special Masses, e.g. Christmas, Ash Wednesday, Easter Sunday,
// Holy Monday, Holy Tuesday, Holy Wednesday, Holy Thursday, Good Friday, Pentecost Sunday, etc.
func NameMassPass(filename, resourceDir string) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// unmarshal calendar as generic map so we can add "name" fields
	calendar := map[string]map[string]any{}
	if err := json.Unmarshal(raw, &calendar); err != nil {
		panic(err)
	}

	// load lectionary resource which contains names for many lectionary keys
	lecRaw, err := os.ReadFile(resourceDir + "/lectionary.json")
	if err != nil {
		panic(err)
	}
	var lec struct {
		Readings map[string]map[string]any `json:"readings"`
	}
	if err := json.Unmarshal(lecRaw, &lec); err != nil {
		panic(err)
	}

	for date, entry := range calendar {
		if entry == nil {
			entry = map[string]any{}
		}

		// preserve existing name if present
		if nameI, ok := entry["name"]; ok {
			if s, ok2 := nameI.(string); ok2 && s != "" {
				calendar[date] = entry
				continue
			}
		}

		lkI, ok := entry["lectionary_key"]
		if !ok {
			calendar[date] = entry
			continue
		}
		lk, ok := lkI.(string)
		if !ok || lk == "" {
			calendar[date] = entry
			continue
		}

		if rd, ok := lec.Readings[lk]; ok && rd != nil {
			if nI, ok := rd["name"]; ok {
				if nStr, ok2 := nI.(string); ok2 && nStr != "" {
					entry["name"] = nStr
				}
			}
		}

		calendar[date] = entry
	}

	out, err := json.MarshalIndent(calendar, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filename, out, 0644); err != nil {
		panic(err)
	}
	fmt.Println(filename)
}
