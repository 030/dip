package docker

import (
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func fileTag(image string) (string, error) {
	b, err := os.ReadFile("Dockerfile")
	if err != nil {
		return "", err
	}
	r, err := regexp.Compile(`FROM ` + image + `:([a-z0-9\.-]+)`)
	if err != nil {
		return "", err
	}
	if !r.Match(b) {
		return "", fmt.Errorf("image: '%s' not found in Dockerfile", image)
	}
	return string(r.FindSubmatch(b)[1]), nil
}

func FileLatest(image, latestTag string) error {
	dft, err := fileTag(image)
	if err != nil {
		return err
	}
	if latestTag != dft {
		return fmt.Errorf("dockerfile tag: '%s' seems to be outdated, as: '%s' exists. Please update the tag in the Dockerfile", dft, latestTag)
	}
	log.Infof("Dockerfile tag: '%s' is up to date. Latest: '%v'", dft, latestTag)
	return err
}

func UpdateFROMStatementDockerfile(image, latestTag string) error {
	log.Infof("Updating FROM statement in the Dockerfile")
	b, err := os.ReadFile("Dockerfile")
	if err != nil {
		return err
	}

	re, err := regexp.Compile(`FROM ` + image + `:([a-z0-9\.-]+)`)
	if err != nil {
		return err
	}
	replaced := re.ReplaceAll(b, []byte(`FROM `+image+`:`+latestTag+`${2}`))
	if err := os.WriteFile("Dockerfile", replaced, 0644); err != nil {
		return err
	}

	return nil
}
