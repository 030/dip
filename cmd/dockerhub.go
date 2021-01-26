package cmd

import (
	"fmt"
	"net/http"
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
		return fmt.Errorf("ResponseCode not 200, but: '%v'. Check whether image: '%v', exists on dockerhub. Perhaps it is an official image and -official is needed", httpStatusCode, image)
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

func LatestDockerHubTagBasedOnRegex(official bool, latest string, image string) string {
	var dockerHubImage string
	if official {
		dockerHubImage = "library/" + image
	} else {
		dockerHubImage = image
	}

	if err := allTags(dockerHubImage, 1); err != nil {
		log.Fatal(err)
	}
	log.Debug(tags)
	log.Debug(len(tags))

	r, err := regexp.Compile(latest)
	if err != nil {
		log.Fatal(err)
	}

	var latestTag string
	for _, tag := range tags {
		if r.MatchString(tag) {
			latestTag = r.FindString(tag)
			break
		}
	}
	tags = tags[:0] // reset slice to prevent that tags related to other image will be found on checking another image
	if latestTag == "" {
		log.Fatal("No tag found. Check whether regex is correct")
	}
	log.Debugf("Latest tag: '%s'", latestTag)
	return latestTag
}
