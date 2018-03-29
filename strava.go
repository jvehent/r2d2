package main

import (
	"fmt"
	"log"
	"strings"
	"time"

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
		log.Println("strava: initiating run")
		activities, err := strava.NewClubsService(client).
			ListActivities(cfg.Strava.ClubID).
			PerPage(50).
			Do()
		if err != nil {
			log.Fatal(err)
		}
		for _, activity := range activities {
			// a shitty cache key, definitely not cryptographically secure
			cachekey := activity.Distance + activity.TotalElevationGain + float64(activity.MovingTime) + float64(activity.ElapsedTime)

			if activityCache.Contains(cachekey) {
				log.Printf("strava: skipping activity from %s %s, already in cache", activity.Athlete.FirstName, activity.Athlete.LastName)
				continue
			}
			if activity.Private {
				log.Printf("strava: skipping private activity %s %s", activity.Athlete.FirstName, activity.Athlete.LastName)
				continue
			}
			activityCache.Add(cachekey, time.Now())
			if isFirstRun {
				log.Printf("strava: skipping activity %s %s because this is my first run and I'm populating my cache", activity.Athlete.FirstName, activity.Athlete.LastName)
				continue
			}
			aDistance := activity.Distance / 1000
			aPace, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(float64(activity.MovingTime)/float64(aDistance))))
			irc.Notice(irchan, fmt.Sprintf("%s %s went for a %0.1f km %s going up %0.1f meters at %s/km.",
				activity.Athlete.FirstName, activity.Athlete.LastName,
				aDistance,
				strings.ToLower(activity.Name),
				activity.TotalElevationGain,
				aPace,
			))
		}
		isFirstRun = false
		time.Sleep(600 * time.Second)
	}
}
