package main

import (
	"fmt"
	"os"

	"github.com/030/dip/internal/docker"
	"github.com/030/dip/internal/k8s"
	"github.com/030/dip/internal/slack"
	"github.com/030/dip/pkg/dockerhub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dockerfile, kubernetes, sendSlackMsg bool
	name, regex                          string
	imageCmd                             = &cobra.Command{
		Use:   "image",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			latestTag, err := dockerhub.LatestTagBasedOnRegex(regex, name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(latestTag)

			channelID := ""
			token := ""
			if sendSlackMsg {
				channelID, err = slackChannelID()
				if err != nil {
					log.Fatal(err)
				}
				token, err = slackToken()
				if err != nil {
					log.Fatal(err)
				}
			}

			if dockerfile {
				if err := docker.FileLatest(name, latestTag); err != nil {
					if sendSlackMsg {
						msg := fmt.Sprintf("Image: '%s' in Dockerfile outdated. Latest tag: '%s'", name, latestTag)
						if os.Getenv("GITLAB_CI") == "true" {
							msg = fmt.Sprintf("%s. CI_PROJECT_PATH: '%s'. BRANCH: '%s'. CI_PROJECT_URL: '%s'", msg, os.Getenv("CI_PROJECT_PATH"), os.Getenv("CI_COMMIT_BRANCH"), os.Getenv("CI_PROJECT_URL"))
						}
						if err := slack.SendMessage(channelID, msg, token); err != nil {
							log.Fatal(err)
						}
					}

					log.Fatal(err)
				}
			}

			if kubernetes {
				images, err := imagesToBeValidated()
				if err != nil {
					log.Fatal(err)
				}

				k := k8s.Images{ToBeValidated: images, SlackToken: token, SlackChannelID: channelID}
				if err := k.UpToDate(); err != nil {
					log.Fatal(err)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Docker image to be checked whether it is up to date")
	if err := imageCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	imageCmd.Flags().StringVarP(&regex, "regex", "r", "", "Regex for finding the latest image tag")
	if err := imageCmd.MarkFlagRequired("regex"); err != nil {
		log.Fatal(err)
	}

	imageCmd.Flags().BoolVar(&dockerfile, "dockerfile", false, "Check whether the image that resides in the Dockerfile is outdated")
	imageCmd.Flags().BoolVar(&kubernetes, "kubernetes", false, "Check whether the image in a k8s file is outdated")
	imageCmd.Flags().BoolVar(&sendSlackMsg, "sendSlackMsg", false, "Send message to Slack")
}
