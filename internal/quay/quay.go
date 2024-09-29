package quay

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/030/dip/internal/sort"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func checkForEachTagWhetherItMatchesTheRegexSortItAndReturnLatestTag(latest string, tags []string) (string, error) {
	r, err := regexp.Compile(latest)
	if err != nil {
		return "", err
	}
	var latestTags []string
	for _, tag := range tags {
		log.Debugf("check whether: '%s', matches regex: '%s'", tag, latest)
		if r.MatchString(tag) {
			latestTags = append(latestTags, r.FindString(tag))
		}
	}
	// tags = tags[:0] // reset slice to prevent that tags related to other image will be found on checking another image
	if latestTags == nil {
		return "", fmt.Errorf("no tags were found. Check whether regex is correct")
	}
	log.Debug(latestTags)
	latestTag, err := sort.Tags(latestTags)
	if err != nil {
		return "", err
	}

	log.Debugf("Latest tag: '%s'", latestTag)
	return latestTag, nil
}

func jsonTags(image string) ([]byte, error) {
	resp, err := http.Get("https://quay.io/api/v1/repository/" + image + "/tag/")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func latest(json []byte, regex string) (string, error) {
	value := gjson.GetBytes(json, "tags.#.name")
	tags := value.Array()
	tagsFromJSON := []string{}
	for _, tag := range tags {
		tagsFromJSON = append(tagsFromJSON, tag.String())
	}

	latestTag, err := checkForEachTagWhetherItMatchesTheRegexSortItAndReturnLatestTag(regex, tagsFromJSON)
	if err != nil {
		return "", err
	}

	return latestTag, nil
}

func LatestTagBasedOnRegex(regex, image string) (string, error) {
	j, err := jsonTags(image)
	if err != nil {
		return "", err
	}

	return latest(j, regex)
}
