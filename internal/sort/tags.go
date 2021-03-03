package sort

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

func semantic(tag string) (bool, error) {
	matched, err := regexp.Match("^([0-9]+\\.){2}[0-9]+", []byte(tag))
	if err != nil {
		return false, err
	}
	return matched, nil
}

func something(tags []string) (string, error) {
	var latestTag string
	var m = make(map[int]string)
	for _, tag := range tags {
		sem, err := semantic(tag)
		if err != nil {
			return "", err
		}
		if !sem {
			return "", fmt.Errorf("Tag: '%s' is not a semantic version", tag)
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

	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for i, k := range keys {
		if i == len(keys)-1 {
			latestTag = m[k]
		}
	}
	return latestTag, nil
}
