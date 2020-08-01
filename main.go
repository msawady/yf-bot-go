package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/slack-go/slack"
)

type Config struct {
	SlackToken string `toml:"slack_token"`
}

func main() {
	var config Config
	var confFileName string
	if value, ok := os.LookupEnv("YF_BOT_CONFIG"); ok {
		confFileName = value
	} else {
		confFileName = "config.toml"
	}
	_, err := toml.DecodeFile(confFileName, &config)
	if err != nil {
		panic(err)
	}
	api := slack.New(
		config.SlackToken,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	usage := `
	Usage: please input stock code split with space after my name.
	example: @yahoo 7201.T 7203.T (max: 5)
	`
	tickerPattern := regexp.MustCompile(`[A-Z]{1,5}`)
	ricPattern := regexp.MustCompile(`\d{4}\.T`)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var myID string
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			myID = ev.Info.User.ID

		case *slack.MessageEvent:
			msg := ev.Text
			if strings.HasPrefix(msg, fmt.Sprintf("<@%s>", myID)) {
				blocks := strings.Split(msg, " ")
				channel := ev.Channel
				messageTs := ev.EventTimestamp
				threadTs := ev.ThreadTimestamp
				fmt.Println("mention at: " + channel)
				fmt.Println("msg: " + msg)
				if len(blocks) < 1 {
					send(rtm, usage, channel, messageTs, threadTs)
				}

				cands := blocks[1:]
				var syms []string
				var invalids []string

				for _, cand := range cands {
					upper := strings.ToUpper(cand)
					if ricPattern.MatchString(upper) || tickerPattern.MatchString(upper) {
						syms = append(syms, upper)
					} else {
						invalids = append(invalids, cand)
					}
				}

				if len(syms) < 1 || len(syms) > 5 {
					send(rtm, usage, channel, messageTs, threadTs)
				} else if len(invalids) > 0 {
					send(rtm, "Following blocks are discarded due to invalid format.\n"+strings.Join(invalids, ","), channel, messageTs, threadTs)
				}

				urlBase := `https://finance.yahoo.com/quote/`
				for _, sym := range syms {
					send(rtm, urlBase+sym, channel, messageTs, threadTs)
				}
			}

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
		}

	}
}

func send(rtm *slack.RTM, message string, channel string, timestamp string, threadTs string) {
	if timestamp == threadTs {
		rtm.SendMessage(rtm.NewOutgoingMessage(message, channel))
	} else {
		rtm.SendMessage(rtm.NewOutgoingMessage(message, channel, withThreadTs(threadTs)))
	}
}

func withThreadTs(ts string) slack.RTMsgOption {
	return func(o *slack.OutgoingMessage) {
		o.ThreadTimestamp = ts
	}
}
