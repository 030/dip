package dockerhub

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
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

func LatestTagBasedOnRegex(official bool, latest string, image string) string {
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
	log.Debugf("Regex: '%s'", latest)

	r, err := regexp.Compile(latest)
	if err != nil {
		log.Fatal(err)
	}

	var latestTags []string
	for _, tag := range tags {
		log.Debugf("Check whether: '%s', matches regex: '%s'", tag, latest)
		if r.MatchString(tag) {
			latestTags = append(latestTags, r.FindString(tag))
		}
	}
	tags = tags[:0] // reset slice to prevent that tags related to other image will be found on checking another image
	if latestTags == nil {
		log.Fatal("No tags were found. Check whether regex is correct")
	}
	log.Debug(latestTags)
	latestTag, err := something(latestTags)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Latest tag: '%s'", latestTag)
	return latestTag
}

func semantic(tag string) (bool, error) {
	matched, err := regexp.Match("^([0-9]+\\.){2}[0-9]+", []byte(tag))
	if err != nil {
		return false, err
	}
	return matched, nil
}

func something(tags []string) (string, error) {
	log.Debugf("Input: '%v'", tags)
	var latestTag string
	var m = make(map[int]string)
	for _, tag := range tags {
		sem, err := semantic(tag)
		if err != nil {
			return "", err
		}
		if !sem {
			continue
			// return "", fmt.Errorf("Tag: '%s' is not a semantic version", tag)
		}

		reg, err := regexp.Compile("[a-zA-Z-]+")
		if err != nil {
			return "", err
		}
		tagWithoutChars := reg.ReplaceAllString(tag, "")
		semver := regexp.MustCompile("\\.").Split(tagWithoutChars, 3)
		major, err := strconv.Atoi(semver[0])
		if err != nil {
			return "", err
		}
		minor, err := strconv.Atoi(semver[1])
		if err != nil {
			return "", err
		}
		patch, err := strconv.Atoi(semver[2])
		if err != nil {
			return "", err
		}
		version := major*1e6 + minor*1e3 + patch
		m[version] = tag
	}
	log.Debugf("Map: '%v'", m)
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	log.Debugf("Sorted keys: '%v'", keys)

	for i, k := range keys {
		if i == len(keys)-1 {
			latestTag = m[k]
		}
	}
	return latestTag, nil
}
