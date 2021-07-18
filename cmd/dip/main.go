package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/030/dip/internal/k8s"
	"github.com/030/dip/pkg/dockerhub"
	sasm "github.com/030/sasm/pkg/slack"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Version string

func dockerfileTag(i string) (string, error) {
	b, err := ioutil.ReadFile("Dockerfile")
	if err != nil {
		return "", err
	}
	r, err := regexp.Compile("FROM " + i + ":(.*)")
	if err != nil {
		return "", err
	}
	if !r.Match(b) {
		return "", fmt.Errorf("no match")
	}
	group := r.FindSubmatch(b)
	log.Debugf("Dockerfile image: '%s' and tag: '%s'", i, string(group[1]))
	return string(group[1]), nil
}

func k8sfileTag(image string) (string, error) {
	b, err := ioutil.ReadFile("deploy.yml")
	if err != nil {
		return "", err
	}
	r, err := regexp.Compile("image: " + image + ":(.*)")
	if err != nil {
		return "", err
	}
	if !r.Match(b) {
		return "", fmt.Errorf("no match")
	}
	group := r.FindSubmatch(b)
	log.Debugf("k8sfile image: '%s' and tag: '%s'", image, string(group[1]))
	return string(group[1]), nil
}

func main() {
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
	dockerfile := flag.Bool("dockerfile", false, "Whether dockerfile should be checked.")
	k8sArg := flag.Bool("k8s", false, "Whether k8s should be checked.")
	slack := flag.Bool("slack", false, "Whether a message should be sent to Slack.")
	image := flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
	official := flag.Bool("official", false, "Use this parameter if an image is official according to dockerhub.")
	latest := flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")
	version := flag.Bool("version", false, "The version of DIP.")

	flag.Parse()

	if (*image == "" || *latest == "") && !*version {
		flag.Usage()
		log.Fatal("image and latest subcommands are mandatory")
	}

	if *version {
		fmt.Println(Version)
		return
	}

	if *debug {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}

	latestTag, err := dockerhub.LatestTagBasedOnRegex(*official, *latest, *image)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(latestTag)

	if *dockerfile {
		dft, err := dockerfileTag(*image)
		if err != nil {
			log.Fatal(err)
		}
		if latestTag != dft {
			msg := fmt.Sprintf("Dockerfile image: '%s' tag: '%s' seems to be outdated, as: '%s' exists. Please update the tag in the Dockerfile", *image, dft, latestTag)
			if *slack {
				viper.SetConfigName("config")
				viper.SetConfigType("yml")

				home, err := homedir.Dir()
				if err != nil {
					log.Fatal(err)
				}
				viper.AddConfigPath(filepath.Join(home, ".dip"))
				if err := viper.ReadInConfig(); err != nil {
					log.Fatalf("Fatal error config file: %w", err)
				}

				log.Info("Sending message to Slack...")
				t := sasm.Text{Type: "mrkdwn", Text: msg}
				b := []sasm.Blocks{{Type: "section", Text: &t}}
				d := sasm.Data{Blocks: b, Channel: "#dip", Icon: ":robot_face:", Username: "dip"}

				slackToken := viper.GetString("slack_token")
				if slackToken == "" {
					log.Fatalf("slack_token should not be empty. Check whether these resides in: '%s'", viper.ConfigFileUsed())
				}
				d.PostMessage(slackToken)
			}
			log.Fatalf(msg)
		} else {
			log.Infof("Dockerfile tag: '%s' is up to date. Latest: '%v'", dft, latestTag)
		}
	}

	if *k8sArg {
		k8s.Images()
		log.Fatal("Doei")

		dft, err := k8sfileTag(*image)
		if err != nil {
			log.Fatal(err)
		}
		if latestTag != dft {
			msg := fmt.Sprintf("Dockerfile image: '%s' tag: '%s' seems to be outdated, as: '%s' exists. Please update the tag in the Dockerfile", *image, dft, latestTag)
			if *slack {
				viper.SetConfigName("config")
				viper.SetConfigType("yml")

				home, err := homedir.Dir()
				if err != nil {
					log.Fatal(err)
				}
				viper.AddConfigPath(filepath.Join(home, ".dip"))
				if err := viper.ReadInConfig(); err != nil {
					log.Fatalf("Fatal error config file: %w", err)
				}

				log.Info("Sending message to Slack...")
				t := sasm.Text{Type: "mrkdwn", Text: msg}
				b := []sasm.Blocks{{Type: "section", Text: &t}}
				d := sasm.Data{Blocks: b, Channel: "#dip", Icon: ":robot_face:", Username: "dip"}

				slackToken := viper.GetString("slack_token")
				if slackToken == "" {
					log.Fatalf("slack_token should not be empty. Check whether these resides in: '%s'", viper.ConfigFileUsed())
				}
				d.PostMessage(slackToken)
			}
			log.Fatalf(msg)
		} else {
			log.Infof("Dockerfile tag: '%s' is up to date. Latest: '%v'", dft, latestTag)
		}
	}
}
