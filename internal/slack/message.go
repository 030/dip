package slack

import (
	"fmt"

	sasm "github.com/030/sasm/pkg/slack"
	log "github.com/sirupsen/logrus"
)

func SendMessage(msg, token string) error {
	if token == "" || msg == "" {
		return fmt.Errorf("slack_token or msg should not be empty")
	}

	log.Info("Sending message to Slack...")
	t := sasm.Text{Type: "mrkdwn", Text: msg}
	b := []sasm.Blocks{{Type: "section", Text: &t}}
	d := sasm.Data{Blocks: b, Channel: "#dip", Icon: ":dip:", Username: "dip"}

	if err := d.PostMessage(token); err != nil {
		return err
	}
	return nil
}
