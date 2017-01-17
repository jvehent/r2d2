package main

import (
	"strings"
	"time"
)

const timeHelp = "print the time at a given timezone, such as 'Europe/Paris', taken from /usr/share/zoneinfo. accepts shortcuts 'poland', 'france', 'sarasota', 'winnipeg', 'pdt'"

func getTimeIn(timezone string) string {
	if timezone != "" {
		switch timezone {
		case "poland":
			timezone = "Europe/Warsaw"
		case "france":
			timezone = "Europe/Paris"
		case "sarasota":
			return "it's always Mojito time in Sarasota!"
		case "winnipeg":
			timezone = "America/Winnipeg"
		case "pdt":
			timezone = "America/Los_Angeles"
		}
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return "invalid time location: " + timezone
		}
		t := time.Now()
		return "the time in " + timezone + " is " + t.In(loc).String()
	}
	return worldtime()
}

func worldtime() (s string) {
	t := time.Now()
	for _, timezone := range []string{"America/Los_Angeles", "America/New_York", "Europe/London", "Europe/Berlin", "Europe/Moscow", "Asia/Taipei", "Australia/Sydney", "Pacific/Auckland"} {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return "invalid time location: " + timezone
		}
		city := strings.Split(timezone, "/")[1]
		switch city {
		case "Berlin":
			city = "Paris/Berlin"
		case "Los_Angeles":
			city = "San Francisco"
		case "New_York":
			city = "New York/Toronto"
		}
		s += city + "=" + t.In(loc).Format("15:04") + " "
	}
	return s
}
