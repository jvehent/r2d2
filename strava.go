package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codingsince1985/geo-golang/google"
	"github.com/hashicorp/golang-lru"
	strava "github.com/strava/go.strava"
	goirc "github.com/thoj/go-ircevent"
)

func stravaHelp() string {
	return fmt.Sprintf("Follow the Strava activities of members of club https://www.strava.com/clubs/%d", cfg.Strava.ClubID)
}

func watchStrava(irc *goirc.Connection) {
	irchan := strings.Split(cfg.Irc.Channels[0], " ")[0]
	if cfg.Strava.Channel != "" {
		irc.Join(cfg.Strava.Channel)
		irchan = strings.Split(cfg.Strava.Channel, " ")[0]
	}
	if cfg.Strava.AccessToken == "" {
		log.Println("strava: missing access token, module disabled")
		return
	}
	if cfg.Strava.ClubID == 0 {
		log.Println("strava: missing club id, module disabled")
		return
	}
	if cfg.Strava.GoogleAPIKey == "" {
		log.Println("strava: missing google geocoding api key, module disabled")
		return
	}
	client := strava.NewClient(cfg.Strava.AccessToken)
	if client == nil {
		log.Println("strava: failed to create client, module disabled")
		return
	}
	activityCache, err := lru.New(2048)
	if err != nil {
		log.Fatal("strava: failed to initialize LRU cache:", err)
	}
	isFirstRun := true
	for {
		activities, err := strava.NewClubsService(client).
			ListActivities(cfg.Strava.ClubID).
			PerPage(50).
			Do()
		if err != nil {
			log.Fatal(err)
		}
		for _, activity := range activities {
			if activityCache.Contains(activity.Id) {
				continue
			}
			activityCache.Add(activity.Id, time.Now())
			if isFirstRun {
				continue
			}
			aDistance := activity.Distance / 1000
			irc.Notice(irchan, fmt.Sprintf("%s %s went for a %0.1fkm %s%s",
				activity.Athlete.FirstName, activity.Athlete.LastName,
				aDistance,
				strings.ToLower(activity.Name),
				))
		}
		isFirstRun = false
		time.Sleep(600 * time.Second)
	}
}
