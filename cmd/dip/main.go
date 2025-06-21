package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/030/dip/internal/docker"
	"github.com/030/dip/internal/k8s"
	"github.com/030/dip/internal/quay"
	"github.com/030/dip/internal/slack"
	"github.com/030/dip/pkg/dockerhub"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const ext = "yml"

var (
	cfgCredHome, Version           string
	debug                          bool
	dockerfile, k8sfile            bool
	kubernetes, quayIo             bool
	sendSlackMsg, updateDockerfile bool
	name, regex                    string
)

var rootCmd = &cobra.Command{
	Use:     "dip",
	Short:   "A Docker/Kubernetes image policy enforcement CLI",
	Long:    "dip helps enforce up-to-date image tags in Dockerfiles and Kubernetes manifests. It also supports Slack notifications and regex-based tag matching from DockerHub and Quay.io.",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetReportCaller(true)
			log.SetLevel(log.DebugLevel)
		}
	},
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Verify and update container image tags",
	Long:  "The image subcommand validates Docker and Kubernetes image tags against latest tags from DockerHub or Quay.io. It can also send Slack notifications if tags are outdated.",
	Run: func(cmd *cobra.Command, args []string) {
		latestTag, err := fetchLatestTag()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(latestTag)

		if dockerfile {
			err = docker.FileLatest(name, latestTag)
			if err != nil {
				sendSlackIfOutdated(latestTag)
				log.Fatal(err)
			}
		}

		if updateDockerfile {
			err = docker.UpdateFROMStatementDockerfile(name, latestTag)
			if err != nil {
				log.Fatal(err)
			}
		}

		if k8sfile {
			tag, err := k8s.FileTag(name)
			if err != nil {
				log.Fatal(err)
			}
			if latestTag != tag {
				log.Fatal(fmt.Errorf("k8sfile tag: '%s' seems to be outdated, as: '%s' exists. Please update the tag in the k8sfile", tag, latestTag))
			}
			log.Infof("k8sfile tag: '%s' is up to date. Latest: '%v'", tag, latestTag)
		}

		if kubernetes {
			images, err := imagesToBeValidated()
			if err != nil {
				log.Fatal(err)
			}
			token, err := slackToken()
			if err != nil {
				log.Fatal(err)
			}
			channelID, err := slackChannelID()
			if err != nil {
				log.Fatal(err)
			}
			k := k8s.Images{ToBeValidated: images, SlackToken: token, SlackChannelID: channelID}
			err = k.UpToDate()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func main() {
	Execute()
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgCredHome, "configCredHome", "", "Config and credential file home directory (default is $HOME/.dip)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	imageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the Docker image to check")
	cobra.CheckErr(imageCmd.MarkFlagRequired("name"))

	imageCmd.Flags().StringVarP(&regex, "regex", "r", "", "Regex to find the latest image tag")
	cobra.CheckErr(imageCmd.MarkFlagRequired("regex"))

	imageCmd.Flags().BoolVar(&dockerfile, "dockerfile", false, "Check if Dockerfile image tag is outdated")
	imageCmd.Flags().BoolVar(&updateDockerfile, "updateDockerfile", false, "Update Dockerfile FROM statement")
	imageCmd.Flags().BoolVar(&k8sfile, "k8sfile", false, "Check if Kubernetes manifest image tags are outdated")
	imageCmd.Flags().BoolVar(&kubernetes, "kubernetes", false, "Run full Kubernetes image tag validation")
	imageCmd.Flags().BoolVar(&quayIo, "quayIo", false, "Check tags on Quay.io instead of DockerHub")
	imageCmd.Flags().BoolVar(&sendSlackMsg, "sendSlackMsg", false, "Send Slack message when outdated image is found")

	rootCmd.AddCommand(imageCmd)
}

func initConfig() {
	log.SetReportCaller(true)
	if debug {
		log.SetLevel(log.DebugLevel)
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
}

func fetchLatestTag() (string, error) {
	if quayIo {
		q := quay.Quay{
			HTTPGetter: quay.HTTPGet{},
			Image:      name,
		}
		return q.LatestTagBasedOnRegex(regex, name)
	}
	return dockerhub.LatestTagBasedOnRegex(regex, name)
}

func sendSlackIfOutdated(latestTag string) {
	if !sendSlackMsg {
		return
	}

	channelID, err := slackChannelID()
	if err != nil {
		log.Fatal(err)
	}

	token, err := slackToken()
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf("Image: '%s' in Dockerfile outdated. Latest tag: '%s'", name, latestTag)
	if os.Getenv("GITLAB_CI") == "true" {
		msg = fmt.Sprintf("%s. CI_PROJECT_PATH: '%s'. BRANCH: '%s'. CI_PROJECT_URL: '%s'",
			msg, os.Getenv("CI_PROJECT_PATH"), os.Getenv("CI_COMMIT_BRANCH"), os.Getenv("CI_PROJECT_URL"))
	}

	err = slack.SendMessage(channelID, msg, token)
	if err != nil {
		log.Fatal(err)
	}
}

func viperBase(filename string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(filepath.Join(home, ".dip"))

	if cfgCredHome != "" {
		viper.AddConfigPath(cfgCredHome)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
	}
	return nil
}

func credsValue(key string) (string, error) {
	err := viperBase("creds")
	if err != nil {
		return "", err
	}
	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("no %s found. Check whether the '%s' variable is populated in '%s'", key, key, viper.ConfigFileUsed())
	}
	return value, nil
}

func slackToken() (string, error) {
	return credsValue("slack_token")
}

func slackChannelID() (string, error) {
	return credsValue("slack_channel_id")
}

func imagesToBeValidated() (map[string]interface{}, error) {
	err := viperBase("config")
	if err != nil {
		return nil, err
	}
	images := viper.GetStringMap("dip_images")
	if len(images) == 0 {
		return nil, fmt.Errorf("no images found. Check whether the 'dip_images' variable is populated in '%s'", viper.ConfigFileUsed())
	}
	return images, nil
}
