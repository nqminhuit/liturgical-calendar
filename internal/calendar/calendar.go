package calendar

import "time"

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
