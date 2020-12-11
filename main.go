package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const dockerRegistry = "https://registry.hub.docker.com/v2/repositories/"

var tags = []string{}

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

func main() {
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
	image := flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
	latest := flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")

	if len(os.Args[1:]) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()

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
	fmt.Println(s)
}
