package icsProcessing

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

type Event struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	TimeStamp   time.Time `json:"timestamp"`
	Summary     string    `json:"summary"`
	Location    string    `json:"location,omitempty"`
	UID         string    `json:"uid"`
	Status      string    `json:"status"`
	Description string    `json:"description,omitempty"`
}

type Calendar struct {
	Name     string  `json:"name"`
	Version  string  `json:"version"`
	Events   []Event `json:"events"`
	TimeZone string  `json:"timezone"`
}

func Parse(icsString string) (Calendar, error) {
	scanner := bufio.NewScanner(strings.NewReader(icsString))
	calendar := Calendar{
		Events: make([]Event, 0),
	}

	var currentEvent *Event

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "BEGIN":
			if value == "VEVENT" {
				currentEvent = &Event{}
			}
		case "END":
			if value == "VEVENT" && currentEvent != nil {
				calendar.Events = append(calendar.Events, *currentEvent)
				currentEvent = nil
			}
		case "X-WR-CALNAME":
			calendar.Name = value
		case "VERSION":
			calendar.Version = value
		case "X-WR-TIMEZONE":
			calendar.TimeZone = value
		}

		if currentEvent != nil {
			switch key {
			case "DTSTART":
				t, err := parseIcsTime(value)
				if err == nil {
					currentEvent.Start = t
				}
			case "DTEND":
				t, err := parseIcsTime(value)
				if err == nil {
					currentEvent.End = t
				}
			case "DTSTAMP":
				t, err := parseIcsTime(value)
				if err == nil {
					currentEvent.TimeStamp = t
				}
			case "SUMMARY":
				currentEvent.Summary = value
			case "UID":
				currentEvent.UID = value
			case "STATUS":
				currentEvent.Status = value
			case "DESCRIPTION":
				currentEvent.Description = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return calendar, fmt.Errorf("błąd podczas skanowania pliku: %w", err)
	}

	return calendar, nil
}

func parseIcsTime(timestamp string) (time.Time, error) {
	// Usuń literę 'Z' z końca, jeśli występuje
	timestamp = strings.TrimSuffix(timestamp, "Z")

	// Format czasu w ICS: YYYYMMDDTHHMMSS
	const icsTimeFormat = "20060102T150405"

	return time.Parse(icsTimeFormat, timestamp)
}
