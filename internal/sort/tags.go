package sort

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func semantic(tag string) (bool, error) {
	matched, err := regexp.Match("^v?([0-9]+\\.){2}[0-9]+", []byte(tag))
	if err != nil {
		return false, err
	}
	return matched, nil
}

func Tags(tags []string) (string, error) {
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
		}

		reg, err := regexp.Compile("[_a-zA-Z-]+")
		if err != nil {
			return "", err
		}
		log.Debugf("tag: '%v'", tag)
		tagWithoutChars := reg.ReplaceAllString(tag, "")
		log.Debugf("tagWithoutChars: '%v'", tagWithoutChars)
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
	if latestTag == "" {
		return "", fmt.Errorf("Cannot find the latest tag. Check whether the tags are semantic")
	}

	return latestTag, nil
}
