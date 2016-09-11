package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	goirc "github.com/thoj/go-ircevent"
)

const githubHelp = "follow commits on multiple github repositories. get the list of followed repos with 'github repos'"

func watchGithub(irc *goirc.Connection) {
	irchan := cfg.Irc.Channel
	if cfg.Github.Channel != "" {
		if cfg.Github.ChannelPass != "" {
			irc.Join(cfg.Github.Channel + " " + cfg.Github.ChannelPass)
		} else {
			irc.Join(cfg.Github.Channel)
		}
		irchan = cfg.Github.Channel
	}
	var err error
	// start the github watcher
	evchan := make(chan string)
	githubCli := makeGithubClient(cfg.Github.Token)
	for _, repo := range cfg.Github.Repos {
		splitted := strings.Split(repo, "/")
		if len(splitted) != 2 {
			irc.Privmsgf(cfg.Irc.Channel, "Invalid repository syntax '%s'. Must be <owner>/<reponame>", repo)
			continue
		}
		// don't run everything at once, we've got time...
		time.Sleep(3 * time.Second)
		go func() {
			for {
				sleepfor := time.Duration(150 + (rand.Int() % 150))
				time.Sleep(sleepfor * time.Second)
				err = followRepoEvents(githubCli, splitted[0], splitted[1], evchan)
				if err != nil {
					log.Println("github follower crashed with error", err)
				}
			}
		}()
	}
	go func() {
		for ev := range evchan {
			// no more than one post per second
			time.Sleep(time.Second)
			irc.Privmsgf(irchan, "%s", ev)
		}
	}()

}

func makeGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func followRepoEvents(cli *github.Client, owner, repo string, evchan chan string) (err error) {
	lastID := "null"
	for {
		events, _, err := cli.Activity.ListRepositoryEvents(owner, repo, nil)
		if err != nil {
			return err
		}
		evctr := 0
		for _, ev := range events {
			evctr++
			// if we're in debug mode, print everything. otherwise, reasons to skip this event:
			// 1. we have already displayed 5 events in this loop
			// 2. the lastid is null, which indicates this is the first loop
			// 3. the current event is the last event, which indicates nothing new happened
			if !cfg.Github.Debug && ((evctr >= 6) || (lastID == "null") || (*ev.ID == lastID)) {
				break
			}
			switch *ev.Type {
			case "PushEvent":
				pe := ev.Payload()
				for _, c := range pe.(*github.PushEvent).Commits {
					evchan <- fmt.Sprintf("\x032[%s/%s]\x03 %s - %s \x038https://github.com/%s/%s/commit/%s\x03",
						owner, repo, *c.Author.Name, *c.Message, owner, repo, *c.SHA)
				}
			}
		}
		if len(events) > 0 {
			lastID = *events[0].ID
		}
		time.Sleep(60 * time.Second)
	}
}

func githubPrintReposList() string {
	return "list of followed github repositories: " + strings.Join(cfg.Github.Repos, ", ")
}
