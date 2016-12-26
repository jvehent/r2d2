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
	return fmt.Sprintf("Follow the Strava activities of members of club https://www.strava.com/dashboard?club_id=%d", cfg.Strava.ClubID)
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
			var aType string
			switch activity.Type.String() {
			case "Run":
				aType = "ran"
			case "Ride":
				aType = "biked"
			case "Hike":
				aType = "hiked"
			case "Kayaking":
				aType = "kayaked"
			default:
				aType = activity.Type.String()
			}
			aDistance := activity.Distance / 1000
			aDuration, err := time.ParseDuration(fmt.Sprintf("%ds", activity.ElapsedTime))
			if err != nil {
				log.Fatal(err)
			}
			aLocation := ""
			geocoder := google.Geocoder(cfg.Strava.GoogleAPIKey)
			address, _ := geocoder.ReverseGeocode(activity.StartLocation[0], activity.StartLocation[1])
			if address != "" {
				addressComp := strings.Split(address, ",")
				if len(addressComp) > 3 {
					aLocation = " around" + strings.Join(addressComp[len(addressComp)-3:], ",")
				} else {
					aLocation = " around " + address
				}
			}
			irc.Privmsg(irchan, fmt.Sprintf("%s %s %s %0.1fkm in %s%s: %s\n",
				activity.Athlete.FirstName, activity.Athlete.LastName,
				aType,
				aDistance,
				aDuration,
				aLocation,
				activity.Name))
		}
		isFirstRun = false
		time.Sleep(60 * time.Second)
	}
}
