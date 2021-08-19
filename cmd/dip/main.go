package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/030/dip/internal/k8s"
	"github.com/030/dip/pkg/dockerhub"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ext = "yml"

var (
	debug, dockerfile, k8sArg, version *bool
	config, image, latest              *string
	Version                            string
)

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

func dockerfileOption() error {
	latestTag, err := dockerhub.LatestTagBasedOnRegex(*latest, *image)
	if err != nil {
		return err
	}
	fmt.Println(latestTag)

	if *dockerfile {
		dft, err := dockerfileTag(*image)
		if err != nil {
			return err
		}
		if latestTag != dft {
			return fmt.Errorf("dockerfile tag: '%s' seems to be outdated, as: '%s' exists. Please update the tag in the Dockerfile", dft, latestTag)
		}
		log.Infof("Dockerfile tag: '%s' is up to date. Latest: '%v'", dft, latestTag)
	}
	return nil
}

func debugOption() {
	if *debug {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
}

func k8sArgOption() error {
	if *k8sArg {
		token, err := slackToken()
		if err != nil {
			return err
		}

		images, err := imagesToBeValidated()
		if err != nil {
			return err
		}

		k := k8s.Images{ToBeValidated: images, SlackToken: token}
		if err := k.UpToDate(); err != nil {
			return err
		}
	}
	return nil
}

func validationOption() error {
	if (*image == "" || *latest == "") && !*version {
		flag.Usage()
		return fmt.Errorf("image and latest subcommands are mandatory")
	}
	return nil
}

func versionOption() {
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
}

func options() error {
	debugOption()

	if err := k8sArgOption(); err != nil {
		return err
	}

	if err := validationOption(); err != nil {
		return err
	}

	versionOption()

	if err := dockerfileOption(); err != nil {
		return err
	}

	return nil
}

func viperBase(path, filename string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(filepath.Join(home, ".dip"))

	if path != "" {
		viper.SetConfigFile(filepath.Join(path, filename+"."+ext))
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
	}
	return nil
}

func slackToken() (string, error) {
	if err := viperBase(*config, "creds"); err != nil {
		return "", err
	}
	slackToken := viper.GetString("slack_token")
	return slackToken, nil
}

func imagesToBeValidated() (map[string]interface{}, error) {
	if err := viperBase(*config, "config"); err != nil {
		return nil, err
	}
	images := viper.GetStringMap("dip_images")
	if len(images) == 0 {
		return nil, fmt.Errorf("no images found. Check whether the 'dip_images' variable is populated in '%s'", viper.ConfigFileUsed())
	}
	return images, nil
}

func main() {
	config = flag.String("config", "", "the file path that contains the configuration")
	debug = flag.Bool("debug", false, "Whether debug mode should be enabled.")
	dockerfile = flag.Bool("dockerfile", false, "Whether dockerfile should be checked.")
	image = flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
	k8sArg = flag.Bool("k8s", false, "Whether images are up to date in a k8s or openshift cluster.")
	latest = flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")
	version = flag.Bool("version", false, "The version of DIP.")

	flag.Parse()

	if err := options(); err != nil {
		log.Fatal(err)
	}
}
