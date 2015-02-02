package main

import (
	"code.google.com/p/gcfg"
	"flag"
	goirc "github.com/thoj/go-ircevent"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Irc struct {
		Server         string
		Channel        string
		Nick, Nickpass string
		TLS            bool
		Debug          bool
	}
	Github struct {
		Token string
		Repos []string
	}
}

var cfg Config

func main() {
	var (
		irc *goirc.Connection
		err error
	)
	var configFile = flag.String("c", "r2d2.cfg", "Load configuration from file")
	flag.Parse()
	_, err = os.Stat(*configFile)
	if err != nil {
		log.Fatal("%v", err)
		os.Exit(1)
	}
	err = gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		log.Fatal("Error in configuration file: %v", err)
		os.Exit(1)
	}
	irc = goirc.IRC(cfg.Irc.Nick, cfg.Irc.Nick)
	irc.UseTLS = cfg.Irc.TLS
	irc.VerboseCallbackHandler = cfg.Irc.Debug
	irc.Debug = cfg.Irc.Debug
	err = irc.Connect(cfg.Irc.Server)
	if err != nil {
		log.Fatal("Connection to IRC server failed: %v", err)
		os.Exit(1)
	}

	// place a callback on nickserv identification and wait until it is done
	if cfg.Irc.Nickpass != "" {
		identwaiter := make(chan bool)
		irc.AddCallback("NOTICE", func(e *goirc.Event) {
			re := regexp.MustCompile("NickServ IDENTIFY")
			if e.Nick == "NickServ" && re.MatchString(e.Message()) {
				irc.Privmsgf("NickServ", "IDENTIFY %s", cfg.Irc.Nickpass)
				identwaiter <- true
			}
		})
		<-identwaiter
		close(identwaiter)
		irc.ClearCallback("NOTICE")
	}
	// we are identified, let's continue
	irc.Join(cfg.Irc.Channel)
	irc.Privmsg(cfg.Irc.Channel, "beep beedibeep dibeep")

	go watchGithub(irc)

	// add callback that captures messages sent to bot
	terminate := make(chan bool)
	irc.AddCallback("PRIVMSG", func(e *goirc.Event) {
		re := regexp.MustCompile("^" + cfg.Irc.Nick + ":(.+)$")
		if re.MatchString(e.Message()) {
			parsed := re.FindStringSubmatch(e.Message())
			if len(parsed) != 2 {
				return
			}
			req := strings.Trim(parsed[1], " ")
			resp := handleRequest(e.Nick, e.Arguments[1], req)
			irc.Privmsgf(cfg.Irc.Channel, "%s: %s", e.Nick, resp)
		}
	})
	<-terminate
	irc.Loop()
	irc.Disconnect()
}

// handleRequest receives a request as a string and attempt to answer it by looking
// at the first word as a keyword.
func handleRequest(nick, srcchan, req string) string {
	command := strings.Split(req, " ")
	switch command[0] {
	case "github":
		if len(command) > 1 && command[1] == "repos" {
			return githubPrintReposList()
		}
		return "try 'help github'"
	case "help":
		if len(command) > 1 {
			return printHelpFor(command[1])
		}
		return "try 'help <command>', supported commands are: time"
	case "time":
		if len(command) > 1 {
			return getTimeIn(command[1])
		}
		return getTimeIn("")
	default:
		return "I do not know how to answer this..."
	}
}

func printHelpFor(command string) string {
	switch command {
	case "github":
		return "follow commits on multiple github repositories. get the list of followed repos with 'github repos'"
	case "time":
		return "print the time at a given timezone, such as 'Europe/Paris', taken from /usr/share/zoneinfo. accepts shortcuts 'poland', 'france', 'sarasota', 'winnipeg', 'pdt'"
	default:
		return "there is no help for " + command
	}
}

func getTimeIn(timezone string) (resp string) {
	if timezone != "" {
		switch timezone {
		case "poland":
			timezone = "Europe/Warsaw"
		case "france":
			timezone = "Europe/Paris"
		case "sarasota":
			timezone = "America/New_York"
		case "winnipeg":
			timezone = "America/Winnipeg"
		case "pdt":
			timezone = "America/Los_Angeles"
		}
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return "invalid time location: " + timezone
		}
		t := time.Now()
		return "the time in " + timezone + " is " + t.In(loc).String()
	}
	resp = time.Now().UTC().String()
	return
}
