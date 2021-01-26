package cmd

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

func LatestDockerHubTagBasedOnRegex(regex *regexp.Regexp, tags []string) string {
	var latestTag string
	for _, tag := range tags {
		if regex.MatchString(tag) {
			latestTag = regex.FindString(tag)
			break
		}
	}
	if latestTag == "" {
		log.Fatal("No tag found. Check whether regex is correct")
	}
	log.Debugf("Latest tag: '%s'", latestTag)
	return latestTag
}
