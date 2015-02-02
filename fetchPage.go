package main

import (
	goirc "github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func fetchPageTitles(irc *goirc.Connection) {
	irc.AddCallback("PRIVMSG", func(e *goirc.Event) {
		rehttp := regexp.MustCompile("(https?://.+)")
		if rehttp.MatchString(e.Message()) {
			url := rehttp.FindStringSubmatch(e.Message())
			if len(url) < 2 {
				return
			}
			title := fetchTitle(url[1])
			log.Printf("Retrieved tile '%s' from url %s\n", title, url[1])
			if title != "" {
				irc.Privmsgf(cfg.Irc.Channel, "Title: %s", title)
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
		log.Println("Get %s returned ", url, resp.Status)
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
		return r[1]
	}
	return ""
}
