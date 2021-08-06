package sort

import (
	"fmt"

	"github.com/030/dip/internal/semantic"
	log "github.com/sirupsen/logrus"
)

func Tags(tags []string) (string, error) {
	log.Debugf("Input: '%v'", tags)

	latestTag, err := semantic.Latest(tags)
	if err != nil {
		return "", err
	}

	if latestTag == "" {
		latestTag = tags[0]
	}

	if latestTag == "" {
		return "", fmt.Errorf("cannot find the latest tag. Check whether the tags are semantic")
	}

	return latestTag, nil
}
