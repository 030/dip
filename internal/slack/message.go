package slack

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func SendMessage(channelID, msg, token string) error {
	if channelID == "" || token == "" || msg == "" {
		return fmt.Errorf("channelID, slack_token or msg should not be empty")
	}

	log.Info("Sending message to Slack...")
	api := slack.New(token)
	channelID, timestamp, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	return nil
}
