package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	goirc "github.com/thoj/go-ircevent"
)

const untappdHelp = "Follow the check-ins of selected users of Untappd.com. Get the list of followed users with 'untappd users'."

type UntappdAPI struct {
	Response UntappdResponse `json:"response"`
}

type UntappdResponse struct {
	Checkins UntappdCheckins `json:"checkins"`
}

type UntappdCheckins struct {
	Count float64       `json:"count"`
	Items []UntappdItem `json:"items"`
}

type UntappdItem struct {
	ID      float64        `json:"checkin_id"`
	Comment string         `json:"checkin_comment"`
	Score   float64        `json:"rating_score"`
	User    UntappdUser    `json:"user"`
	Beer    UntappdBeer    `json:"beer"`
	Brewery UntappdBrewery `json:"brewery"`
}

type UntappdUser struct {
	Name string `json:"user_name"`
}
type UntappdBeer struct {
	Name  string  `json:"beer_name"`
	Style string  `json:"beer_style"`
	Abv   float64 `json:"beer_abv"`
}

type UntappdBrewery struct {
	Name    string `json:"brewery_name"`
	Country string `json:"country_name"`
}

func watchUntappd(irc *goirc.Connection) {
	var (
		userEvents []string
		err        error
	)
	irchan := cfg.Irc.Channel
	if cfg.Untappd.Channel != "" {
		if cfg.Untappd.ChannelPass != "" {
			irc.Join(cfg.Untappd.Channel + " " + cfg.Untappd.ChannelPass)
		} else {
			irc.Join(cfg.Untappd.Channel)
		}
		irchan = cfg.Untappd.Channel
	}
	lastCheckins := make(map[string]float64)
	for {
		for _, user := range cfg.Untappd.Users {
			sleepfor := time.Duration(30 + (rand.Int() % 30))
			time.Sleep(sleepfor * time.Second)
			// store the last checkin ID, to avoid printing the same checkin twice
			var lastcheckin float64
			if _, ok := lastCheckins[user]; ok {
				lastcheckin = lastCheckins[user]
			}
			if cfg.Untappd.Debug {
				fmt.Println("query untappd activity for", user, "with last checkin set to", lastcheckin)
			}
			userEvents, lastCheckins[user], err = getUntappdActivityFor(user, lastcheckin)
			if err != nil {
				log.Println("Failed to get", user, "'s Untappd activity:", err)
			} else {
				for _, ev := range userEvents {
					irc.Privmsgf(irchan, "%s", ev)
				}
			}
		}
	}
}

func getUntappdActivityFor(user string, lastcheckin float64) (userEvents []string, newlastcheckin float64, err error) {
	target := fmt.Sprintf("https://api.untappd.com/v4/user/checkins/%s?client_id=%s&client_secret=%s",
		user, cfg.Untappd.ClientID, cfg.Untappd.ClientSecret)
	if cfg.Untappd.Debug {
		fmt.Println("querying untappd api at", target)
	}
	resp, err := http.Get(target)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Untappd API returned error " + resp.Status)
		return
	}
	var r UntappdAPI
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Failed to read response from Untappd API")
		return
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		err = fmt.Errorf("Invalid response from Untappd API")
		fmt.Printf("%s\n", body)
		return
	}
	if cfg.Untappd.Debug {
		fmt.Println("retrieved", r.Response.Checkins.Count, "items from untappd api")
	}
	for _, item := range r.Response.Checkins.Items {
		// if this is the first run, don't return any event, just capture the id
		// of the last check in
		if lastcheckin == 0 {
			if cfg.Untappd.Debug {
				log.Println("first run, setting user's last checkin to", item.ID)
			}
			lastcheckin = item.ID
			newlastcheckin = item.ID
			return
		}
		if lastcheckin == item.ID {
			// this item and the followings have already been seen
			break
		}
		userEvents = append(userEvents, fmt.Sprintf("%s is drinking a %s; %.1f%% %s from %s, %s. Score: %.1f/5 - %s\n",
			item.User.Name, item.Beer.Name, item.Beer.Abv, item.Beer.Style, item.Brewery.Name,
			item.Brewery.Country, item.Score, item.Comment))
	}
	// store the ID of the first item as the last checked in
	if len(r.Response.Checkins.Items) > 0 {
		newlastcheckin = r.Response.Checkins.Items[0].ID
	}
	return
}

func untappdPrintUsers() string {
	return "list of followed untappd users: " + strings.Join(cfg.Untappd.Users, ", ")
}
