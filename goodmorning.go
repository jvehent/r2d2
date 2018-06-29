package main

import (
	"math/rand"
	"regexp"
	"time"

	goirc "github.com/thoj/go-ircevent"
)

var goodMorningRe = regexp.MustCompile("(?i)(morning|ohai|howdy|what's up|hiya|hey)")

func sayGoodMorning(irc *goirc.Connection) {
	for {
		if time.Now().UTC().Hour() == cfg.Morning.Hour {
			if time.Now().Weekday() > 0 && time.Now().Weekday() < 6 {
				// only during weekdays
				irc.Privmsgf(cfg.Morning.Channel, "Good Morning %s", cfg.Morning.Who)
			}
		}
		time.Sleep(60*time.Minute + 37*time.Second)
	}
}

func goodMorning() string {
	var mornings = []string{
		"good morning to you",
		"howdy",
		"and a pleasant day to you as well",
		"how's it going?",
		"good to see you!",
		"what's up?",
		"yo!",
		"how've you been?",
		"how do you do?",
		"alright mate?",
		"hiya!",
		"sup?",
		"morning, sunshine!",
		"happy morning",
		"namaste",
		"hakuna matata",
		"top o’ the mornin’!",
		"salut!",
		"buna dimineata",
		"bien le bonjour!",
		"guten morgen",
		"mornin' y'all",
	}
	rand.Seed(time.Now().Unix())
	return mornings[rand.Intn(len(mornings))]
}
