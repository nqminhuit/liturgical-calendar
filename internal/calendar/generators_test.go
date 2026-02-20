package calendar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateYear(t *testing.T) {
	calendar := GenerateYear(2024)
	if len(calendar) != 366 { // 2024 leap year
		t.Errorf("GenerateYear(2024) length = %d, want 366", len(calendar))
	}
	// Check a specific date
	day, ok := calendar["2024-01-01"]
	if !ok {
		t.Error("2024-01-01 not found")
	}
	if day.Season != SeasonChristmas {
		t.Errorf("2024-01-01 season = %s, want %s", day.Season, SeasonChristmas)
	}
	if day.SundayCycle != "B" {
		t.Errorf("2024-01-01 sunday cycle = %s, want B", day.SundayCycle)
	}
	if day.WeekdayCycle != "II" {
		t.Errorf("2024-01-01 weekday cycle = %s, want II", day.WeekdayCycle)
	}
	if day.Weekday != "mon" {
		t.Errorf("2024-01-01 weekday = %s, want mon", day.Weekday)
	}
}

func TestWeekPass(t *testing.T) {
	calendar := GenerateYear(2025)
	data, err := json.Marshal(calendar)
	if err != nil {
		t.Fatal(err)
	}

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	WeekPass(tempFile)

	resultData, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	result := map[string]DayInfo{}
	if err := json.Unmarshal(resultData, &result); err != nil {
		t.Fatal(err)
	}

	// Check specific dates
	tests := map[string]int{
		// Previous Christmas
		"2025-01-01": 1,
		"2025-01-02": 1,
		"2025-01-03": 1,
		"2025-01-04": 1,
		"2025-01-05": 2,
		"2025-01-06": 2,
		"2025-01-07": 2,
		"2025-01-11": 2,
		"2025-01-12": 3,

		// Ordinary time: after Christmas, until Ash Wednesday Mar 5
		"2025-01-13": 1,
		"2025-01-14": 1,
		"2025-01-15": 1,
		"2025-01-16": 1,
		"2025-01-17": 1,
		"2025-01-18": 1,
		"2025-01-19": 2,
		"2025-01-20": 2,
		"2025-02-12": 5,
		"2025-02-23": 7,
		"2025-02-28": 7,
		"2025-03-01": 7,
		"2025-03-02": 8,
		"2025-03-03": 8,
		"2025-03-04": 8,

		// Lent 2025: Ash Wednesday Feb 26, until Holy Saturday Apr 19
		"2025-03-05": 0,
		"2025-03-06": 0,
		"2025-03-07": 0,
		"2025-03-08": 0,
		"2025-03-09": 1,
		"2025-03-10": 1,
		"2025-04-11": 5,
		"2025-04-17": 6,
		"2025-04-18": 6,
		"2025-04-19": 6,

		// Easter 2025: Apr 20
		"2025-04-20": 1,
		"2025-04-21": 1,
		"2025-04-22": 1,
		"2025-04-23": 1,
		"2025-04-24": 1,
		"2025-04-25": 1,
		"2025-04-26": 1,
		"2025-04-27": 2,
		"2025-04-28": 2,
		"2025-04-29": 2,
		"2025-04-30": 2,
		"2025-05-04": 3,
		"2025-05-18": 5,
		"2025-05-19": 5,
		"2025-05-20": 5,
		"2025-05-21": 5,
		"2025-05-22": 5,
		"2025-05-23": 5,
		"2025-05-24": 5,
		"2025-05-25": 6,
		"2025-05-26": 6,
		"2025-05-27": 6,
		"2025-06-01": 7,
		"2025-06-02": 7,
		"2025-06-03": 7,
		"2025-06-04": 7,
		"2025-06-05": 7,
		"2025-06-06": 7,
		"2025-06-07": 7,
		"2025-06-08": 8,

		// Ordinary time: after Pentecost Jun 8
		"2025-06-09": 10,
		"2025-06-10": 10,
		"2025-06-11": 10,
		"2025-06-12": 10,
		"2025-06-13": 10,
		"2025-06-14": 10,
		"2025-06-15": 11,
		"2025-06-16": 11,
		"2025-06-17": 11,
		"2025-06-18": 11,
		"2025-06-19": 11,
		"2025-06-20": 11,
		"2025-06-21": 11,
		"2025-06-22": 12,
		"2025-06-23": 12,
		"2025-06-24": 12,
		"2025-06-25": 12,
		"2025-06-26": 12,
		"2025-06-27": 12,
		"2025-06-28": 12,
		"2025-06-29": 13,
		"2025-06-30": 13,
		"2025-07-20": 16,
		"2025-07-30": 17,
		"2025-08-10": 19,
		"2025-08-20": 20,
		"2025-09-10": 23,
		"2025-09-30": 26,
		"2025-10-20": 29,
		"2025-10-30": 30,
		"2025-11-10": 32,

		// Advent
		"2025-11-30": 1,
		"2025-12-01": 1,
		"2025-12-02": 1,
		"2025-12-03": 1,
		"2025-12-04": 1,
		"2025-12-05": 1,
		"2025-12-06": 1,
		"2025-12-07": 2,
		"2025-12-13": 2,
		"2025-12-14": 3,
		"2025-12-15": 3,
		"2025-12-20": 3,
		"2025-12-21": 4,
		"2025-12-22": 4,
		"2025-12-23": 4,
		"2025-12-24": 4,

		// Christmas
		"2025-12-25": 0,
		"2025-12-26": 0,
		"2025-12-27": 0,
		"2025-12-28": 1,
		"2025-12-29": 1,
		"2025-12-30": 1,
		"2025-12-31": 1,

		// Ordinary after Christmas.
	}

	for date, expected := range tests {
		day, ok := result[date]
		if !ok {
			t.Errorf("Date %s not found", date)
			continue
		}
		if day.WeekOfSeason != expected {
			t.Errorf("Date %s week_of_season = %d, want %d", date, day.WeekOfSeason, expected)
		}
	}
}
