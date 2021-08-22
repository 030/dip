package semantic

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Check verifies whether a version is a semantic version
func Check(tags []string) ([]string, error) {
	var semanticTags []string
	for _, tag := range tags {
		matched, err := regexp.Match("^v?([0-9]+\\.){2}[0-9]+", []byte(tag))
		if err != nil {
			return nil, err
		}
		if matched {
			semanticTags = append(semanticTags, tag)
		}
	}
	return semanticTags, nil
}

// toInt makes it possible to sort semnatic versions, i.e. the highest integer
// represents the latest version
func toInt(tagWithoutChars string) (int, error) {
	semver := regexp.MustCompile(`\.`).Split(tagWithoutChars, 3)
	major, err := strconv.Atoi(semver[0])
	if err != nil {
		return 0, err
	}
	minor, err := strconv.Atoi(semver[1])
	if err != nil {
		return 0, err
	}
	patch, err := strconv.Atoi(semver[2])
	if err != nil {
		return 0, err
	}
	version := major*1e6 + minor*1e3 + patch
	return version, nil
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

func tagWithoutChars(semanticTag string) (int, error) {
	log.Debugf("tag: '%v'", semanticTag)

	re := regexp.MustCompile(`((\d+\.){2}\d+)`)
	match := re.FindStringSubmatch(semanticTag)
	if len(match) == 0 {
		return 0, fmt.Errorf("no match was found. Verify whether the semanticTag: '%v' contains characters", semanticTag)
	}
	tagWithoutChars := match[1]
	log.Debugf("tagWithoutChars: '%v'", tagWithoutChars)

	version, err := toInt(tagWithoutChars)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func sortAndGetLatestTag(m map[int]string) string {
	latestTag := ""
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
	return latestTag
}

func boo(semanticTags []string) (string, error) {
	m := make(map[int]string)
	for _, semanticTag := range semanticTags {
		version, err := tagWithoutChars(semanticTag)
		if err != nil {
			return "", err
		}
		if val, ok := m[version]; ok {
			log.Debugf("Version: '%d' with Tag: '%s' resides already in map!", version, val)

			u1, err := update(val)
			if err != nil {
				return "", err
			}
			log.Debug(u1)

			u2, err := update(semanticTag)
			if err != nil {
				return "", err
			}
			log.Debug(u1)
			if u1 == 0 || u2 == 0 {
				continue
			}
			if u1 > u2 {
				semanticTag = val
			}
			log.Debug(semanticTag)
		}

		m[version] = semanticTag
	}
	latestTag := sortAndGetLatestTag(m)
	return latestTag, nil
}

func Latest(tags []string) (string, error) {
	if len(tags) == 0 {
		return "", fmt.Errorf("tags should not be empty")
	}

	semanticTags, err := Check(tags)
	if err != nil {
		return "", err
	}
	if len(semanticTags) == 0 {
		log.Debug("No semantic tags found")
	}
	return boo(semanticTags)
}
