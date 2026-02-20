# Liturgical Calendar Generator

A Go program that generates liturgical calendar data for any given year, including liturgical seasons, Sunday and weekday cycles, and week numbers within each season.

## Features

- Calculates liturgical seasons (Advent, Christmas, Lent, Easter, Ordinary Time)
- Determines Sunday cycle (A, B, C) based on the 3-year rotation
- Sets weekday cycle (I or II) based on even/odd years
- Computes week numbers within each liturgical season
- Generates JSON output for easy integration

## Installation

1. Ensure you have Go installed (version 1.26 or later recommended)
2. Clone or download this repository

## Usage

Run the program with a year as argument:

```bash
go run main.go <year>
```

For example:
```bash
go run main.go 2025
```

This will generate two JSON files:
- `liturgical-calendar-2025.json`: Basic calendar data
- `liturgical-calendar-2025-weekpass.json`: Calendar with week numbers added

## Output Format

Each day entry in the JSON contains:

```json
{
  "2026-01-01": {
    "season": "christmas",
    "sunday_cycle": "A",
    "weekday_cycle": "II",
    "weekday": "thu",
    "week_of_season": 1
  }
}
```

## Resources

- `resources/lectionary.json`: Contains lectionary readings (to be integrated)
- `resources/sample/`: Sample generated calendars for 2025 and 2026

## Algorithm Notes

- Liturgical year starts on the First Sunday of Advent
- Easter date calculated using Meeus/Jones/Butcher algorithm
- Sunday cycles repeat every 3 years, starting with cycle A in 2008
- Weekday cycles alternate yearly (even years: II, odd years: I)

## Dependencies

- Standard Go library only (time, json, os packages)
