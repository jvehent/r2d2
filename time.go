package main

import (
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
	return time.Now().UTC().String()
}
