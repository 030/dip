package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const dockerRegistry = "https://registry.hub.docker.com/v2/repositories/"

var tags = []string{}
var ver string

func allTags(image string, page int) error {
	resp, err := grequests.Get(dockerRegistry+image+"/tags?page="+strconv.Itoa(page)+"&page_size=100", nil)
	if err != nil {
		return err
	}
	httpStatusCode := resp.StatusCode
	if httpStatusCode != http.StatusOK {
		return fmt.Errorf("ResponseCode not 200, but: '%v'. Check whether image: '%v', exists on dockerhub", httpStatusCode, image)
	}

	tags = append(tags, tagFromJSON(resp.Bytes())...)
	if gjson.GetBytes(resp.Bytes(), "next").String() != "" {
		log.Debug(gjson.GetBytes(resp.Bytes(), "next"))
		log.Debug(page)
		page++
		if err := allTags(image, page); err != nil {
			return err
		}
	}
	return nil
}

func tagFromJSON(b []byte) []string {
	tags := gjson.GetBytes(b, "results.#.name").Array()
	tagsFromJSON := []string{}
	for _, tag := range tags {
		tagsFromJSON = append(tagsFromJSON, tag.String())
	}
	return tagsFromJSON
}

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
		return "", fmt.Errorf("No match")
	}
	group := r.FindSubmatch(b)
	log.Debugf("Dockerfile image: '%s' and tag: '%s'", i, string(group[1]))
	return string(group[1]), nil
}

func main() {
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
	dockerfile := flag.Bool("dockerfile", false, "Whether dockerfile should be checked.")
	image := flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
	latest := flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")
	version := flag.Bool("version", false, "The version of DIP.")

	flag.Parse()

	if (*image == "" || *latest == "") && !*version {
		flag.Usage()
		log.Fatal("image and latest subcommands are mandatory")
	}

	if *version {
		fmt.Println(ver)
		return
	}

	log.SetReportCaller(true)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := allTags(*image, 1); err != nil {
		log.Fatal(err)
	}
	log.Debug(tags)
	log.Debug(len(tags))

	r, err := regexp.Compile(*latest)
	if err != nil {
		log.Fatal(err)
	}

	var s string
	for _, tag := range tags {
		if r.MatchString(tag) {
			s = r.FindString(tag)
			break
		}
	}
	if s == "" {
		log.Fatal("No tag found. Check whether regex is correct")
	}
	log.Infof("Latest tag: '%s' found for image: '%s'", s, *image)

	if *dockerfile {
		dft, err := dockerfileTag(*image)
		if err != nil {
			log.Fatal(err)
		}
		if s != dft {
			log.Fatalf("Dockerfile tag: '%s' seems to be outdated: '%s'", dft, s)
		} else {
			log.Infof("Dockerfile tag: '%s' is up to date. Latest: '%v'", dft, s)
		}
	}
}
