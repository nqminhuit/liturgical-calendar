package calendar

import "testing"

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
