# Liturgical Calendar Generator

A Go program that generates liturgical calendar data for any given year, including liturgical seasons, Sunday and weekday cycles, and week numbers within each season.

## Project Structure

```
liturgical-calendar/
├── cmd/
│   └── main.go              # Command-line entry point
├── internal/
│   └── calendar/
│       ├── calendar.go      # Core calendar calculations (Easter, seasons, cycles)
│       ├── calendar_test.go # Tests for calendar calculations
│       ├── generators.go    # Data generation functions
│       ├── generators_test.go # Tests for generation functions
│       └── types.go         # Data structures and constants
├── resources/
│   └── lectionary.json      # Lectionary readings (not currently integrated)
├── go.mod
├── go.sum
├── Makefile                 # Build and test automation
└── README.md
```

## Features

- Calculates liturgical seasons (Advent, Christmas, Lent, Easter, Ordinary Time)
- Determines Sunday cycle (A, B, C) based on the 3-year rotation
- Sets weekday cycle (I or II) based on even/odd years
- Computes week numbers within each liturgical season
- Generates lectionary keys for each day to facilitate integration with lectionary readings
- Generates JSON output for easy integration with lectionary
- Modular design for maintainability and testing

## Installation

1. Ensure you have Go installed (version 1.26 or later recommended)
2. Clone or download this repository

## Usage

Run the program with a year as argument:

```bash
make build   # Build the binary
make test    # Run tests
make run YEAR=2025  # Run with a specific year
```

For example:
```bash
make run YEAR=2025
```

This will generate a JSON file at `resources/liturgical-calendar-2025.json` containing the complete calendar with seasons, cycles, week numbers, and lectionary keys for integration with lectionary readings.

## Building and Testing

Use the provided Makefile for development tasks:

- `make build`: Compile the binary to `bin/lc`
- `make test`: Run all unit tests with verbose output
- `make compile`: Check if code compiles without building binary
- `make clean`: Remove built binaries
- `make fmt`: Format Go code
- `make vet`: Run Go vet for code analysis

## Output Format

Each day entry in the JSON contains:

```json
{
  "2026-01-01": {
    "season": "christmas",
    "sunday_cycle": "A",
    "weekday_cycle": "II",
    "weekday": "thu",
    "week_of_season": 1,
    "lectionary_key": "christmas_1_thu_II"
  }
}
```

## Resources

- `resources/lectionary.json`: Contains lectionary readings (to be integrated)

## Algorithm Notes

- Liturgical year starts on the First Sunday of Advent
- Easter date calculated using Meeus/Jones/Butcher algorithm
- Sunday cycles repeat every 3 years, starting with cycle A in 2008
- Weekday cycles alternate yearly (even years: II, odd years: I)
- Ordinary Time week numbering: always 34 weeks total, with the first segment from after Christmas to Lent (weeks 1 onward) and the second segment backward from Advent to Pentecost (weeks 34 downward)


