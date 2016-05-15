package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	goirc "github.com/thoj/go-ircevent"
)

func fetchPageTitles(irc *goirc.Connection) {
	irc.AddCallback("PRIVMSG", func(e *goirc.Event) {
		rehttp := regexp.MustCompile("(https?://.+)")
		if rehttp.MatchString(e.Message()) {
			url := rehttp.FindStringSubmatch(e.Message())
			if len(url) < 2 {
				log.Printf("Could not find a message body to work with. event=%+V", e)
				return
			}
			irchan := cfg.Irc.Channel
			if len(e.Arguments) > 0 {
				irchan = e.Arguments[0]
			}
			title := fetchTitle(url[1])
			log.Printf("Retrieved tile '%s' from url %s\n", title, url[1])
			if title != "" {
				irc.Privmsgf(irchan, "Title: %s", title)
			}
		}
	})
	return
}
func fetchTitle(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to retrieve URL", url)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Get %q returned %q", url, resp.Status)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Failed to read response from target"
	}
	re := regexp.MustCompile("<title>(.+)</title>")
	if re.Match(body) {
		r := re.FindStringSubmatch(string(body))
		if len(r) < 2 {
			return ""
		}
		// convert some common html escape sequences back to readable strings
		title := r[1]
		title = strings.Replace(title, "&ndash;", "-", -1)
		title = strings.Replace(title, "&quot;", "\"", -1)
		title = strings.Replace(title, "&#39;", "'", -1)
		title = strings.Replace(title, "&#10;", " ", -1)
		return title
	}
	return ""
}
