package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New(
		"MY_TOKEN",
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	usage := `
	Usage: please input stock code split with space after my name.
	example: @Yahoo Finance 7201.T 7203.T (max: 5)
	`
	ricPattern := regexp.MustCompile(`\d{4}\.T`)
	mentionToMe := "MY_ID"
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			msg := ev.Text
			if strings.HasPrefix(msg, mentionToMe) {
				blocks := strings.Split(msg, " ")
				channel := ev.Channel
				messageTs := ev.EventTimestamp
				threadTs := ev.ThreadTimestamp
				if len(blocks) < 1 {
					send(rtm, usage, channel, messageTs, threadTs)
				}

				cands := blocks[1:len(blocks)]
				var rics []string
				var invalids []string

				for _, b := range cands {
					if ricPattern.MatchString(b) {
						rics = append(rics, b)
					} else {
						invalids = append(invalids, b)
					}
				}

				if len(rics) < 1 || len(rics) > 5 {
					send(rtm, usage, channel, messageTs, threadTs)
				} else if len(invalids) > 0 {
					send(rtm, "Following blocks are discarded due to invalid format.\n"+strings.Join(invalids, ","), channel, messageTs, threadTs)
				}

				urlBase := `https://stocks.finance.yahoo.co.jp/stocks/detail/?code=`
				for _, ric := range rics {
					send(rtm, urlBase+ric, channel, messageTs, threadTs)
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