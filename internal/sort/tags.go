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

		log.Debugf("tag: '%v'", tag)

		re := regexp.MustCompile(`((\d+\.){2}\d+)`)
		match := re.FindStringSubmatch(tag)
		tagWithoutChars := match[1]
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

		if val, ok := m[version]; ok {
			log.Debugf("Version: '%d' with Tag: '%s' resides already in map!", version, val)

			u1, err := update(val)
			if err != nil {
				return "", err
			}
			log.Debug(u1)

			u2, err := update(tag)
			if err != nil {
				return "", err
			}
			log.Debug(u1)
			if u1 == 0 || u2 == 0 {
				continue
			}
			if u1 > u2 {
				tag = val
			}
			log.Debug(tag)
		}

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

func update(tag string) (int, error) {
	re := regexp.MustCompile(`(_|-\w+\d\.)(\d+)`)
	match := re.FindStringSubmatch(tag)
	elements := len(match)
	if elements != 2 {
		log.Debugf("Cannot determine update number. Required elements is two, but was: '%d'", elements)
		return 0, nil
	}
	updateString := match[2]
	u, err := strconv.Atoi(updateString)
	if err != nil {
		return 0, err
	}
	return u, nil
}
