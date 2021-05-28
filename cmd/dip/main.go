package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/030/dip/pkg/dockerhub"
	log "github.com/sirupsen/logrus"
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

func main() {
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
	dockerfile := flag.Bool("dockerfile", false, "Whether dockerfile should be checked.")
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
			log.Fatalf("Dockerfile tag: '%s' seems to be outdated, as: '%s' exists. Please update the tag in the Dockerfile", dft, latestTag)
		} else {
			log.Infof("Dockerfile tag: '%s' is up to date. Latest: '%v'", dft, latestTag)
		}
	}
}
