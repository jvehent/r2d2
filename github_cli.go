package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
	goirc "github.com/thoj/go-ircevent"
	"log"
	"strings"
	"time"
)

const githubHelp = "follow commits on multiple github repositories. get the list of followed repos with 'github repos'"

func watchGithub(irc *goirc.Connection) {
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
		time.Sleep(time.Second)
		go func() {
			for {
				err = followRepoEvents(githubCli, splitted[0], splitted[1], evchan)
				if err != nil {
					log.Println("github follower crashed with error", err)
				}
				time.Sleep(60 * time.Second)
			}
		}()
	}
	go func() {
		for ev := range evchan {
			// no more than one post per second
			time.Sleep(time.Second)
			irc.Privmsgf(cfg.Irc.Channel, "%s", ev)
		}
	}()

}

func makeGithubClient(token string) *github.Client {
	if token != "" {
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token},
		}
		return github.NewClient(t.Client())
	}
	return github.NewClient(nil)
}

func followRepoEvents(cli *github.Client, owner, repo string, evchan chan string) (err error) {
	lastID := "null"
	for {
		events, _, err := cli.Activity.ListRepositoryEvents(owner, repo, nil)
		if err != nil {
			return err
		}
		for _, ev := range events {
			// unless we're in debug mode, we don't print past events
			if !cfg.Github.Debug && lastID == "null" {
				break
			}
			if *ev.ID == lastID {
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
	list := "list of followed github repositories: "
	for _, repo := range cfg.Github.Repos {
		list += repo + ", "
	}
	return list
}
