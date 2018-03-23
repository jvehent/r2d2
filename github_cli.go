package main

import (
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
	irchan := strings.Split(cfg.Irc.Channels[0], " ")[0]
	if cfg.Github.Channel != "" {
		irc.Join(cfg.Github.Channel)
		irchan = cfg.Github.Channel
	}
	var err error
	// start the github watcher
	githubCli := makeGithubClient(cfg.Github.Token)
	for _, repo := range cfg.Github.Repos {
		destinationChannel := irchan
		// first split on whitespaces. first part is the repo,
		// optional second and third are the irc chan and pass
		splitted := strings.Split(repo, " ")
		reposplit := strings.Split(splitted[0], "/")
		if len(reposplit) != 2 {
			irc.Privmsgf(irchan, "Invalid repository syntax '%s'. Must be 'owner/reponame <optional:#ircchan> <optional:channelpass>'", repo)
			continue
		}
		repoOwner := reposplit[0]
		repoName := reposplit[1]
		// if there's a custom channel to send messages to, store it
		if len(splitted) > 1 {
			irc.Join(strings.Join(splitted[1:], " "))
			destinationChannel = splitted[1]
		}
		// don't run everything at once, we've got time...
		time.Sleep(3 * time.Second)
		go func() {
			for {
				log.Printf("github: following repository %s/%s and sending messages into %s", repoOwner, repoName, destinationChannel)
				err = followRepoEvents(irc, githubCli, repoOwner, repoName, destinationChannel)
				if err != nil {
					log.Println("github follower crashed with error", err)
				}
				sleepfor := time.Duration(150 + (rand.Int() % 150))
				time.Sleep(sleepfor * time.Second)
			}
		}()
	}

}

func makeGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func followRepoEvents(irc *goirc.Connection, cli *github.Client, owner, repo, irchan string) (err error) {
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
				if !strings.HasSuffix(*pe.(*github.PushEvent).Ref, "/master") {
					// we only care about commits on the master branch
					if cfg.Github.Debug {
						log.Println("github: ignoring non-master change on ref", *pe.(*github.PushEvent).Ref)
					}
					continue
				}
				for _, c := range pe.(*github.PushEvent).Commits {
					if strings.Contains(*c.Message, "Merge pull request #") {
						continue
					}
					time.Sleep(time.Second)
					msg := strings.Replace(strings.Replace(*c.Message, "\n", " ", -1), "\r", " ", -1)
					irc.Privmsgf(irchan, "\x032[%s/%s]\x03 %s - %s \x038https://github.com/%s/%s/commit/%s\x03",
						owner, repo, *c.Author.Name, msg, owner, repo, *c.SHA)
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
	str := "list of followed github repositories: "
	for _, repo := range cfg.Github.Repos {
		str += strings.Split(repo, " ")[0] + ","
	}
	return str
}
