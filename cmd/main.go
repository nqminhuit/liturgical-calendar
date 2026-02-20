package main

import (
	"encoding/json"
	"fmt"
	"liturgical-calendar/internal/calendar"
	"os"
	"strconv"
)

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
	cal := calendar.GenerateYear(year)

	b, err := json.MarshalIndent(cal, "", "  ")
	if err != nil {
		panic(err)
	}
	filename := fmt.Sprintf("liturgical-calendar-%d.json", year)
	os.WriteFile(filename, b, 0644)
	calendar.WeekPass(filename)
}
