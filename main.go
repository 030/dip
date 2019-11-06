package main

import (
	"flag"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func command(s string) error {
	out, err := exec.Command("bash", "-c", s).CombinedOutput()
	outString := string(out)

	log.Debug(s)
	log.Info(outString)

	if err != nil {
		return err
	}
	return nil
}

// absent checks whether a specific docker image is absent
// in a given docker registry. If true, then a certain docker
// image is absent in a certain docker registry
func absent(image, registry string) bool {
	cmd := "docker pull " + registry + image
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	outString := string(out)

	log.WithFields(log.Fields{
		"cmd":      cmd,
		"output":   outString,
		"image":    image,
		"registry": registry,
	}).Debug("Whether an image is absent in a registry")

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitError.ExitCode()
			log.Debug(exitError.ExitCode())
			return true
		}
	}

	return false
}

func main() {
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled")
	image := flag.String("image", "", "The origin of the image, e.g. nginx:1.17.5-alpine")
	registry := flag.String("registry", "", "To what destination the image should be transferred, e.g. quay.io/some-org/. Note: do not omit the last forward slash.")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	dockerImageAbsent := absent(*image, *registry)

	if !dockerImageAbsent {
		log.Fatal("Docker image: ", *image, " already exists in registry: ", *registry)
	} else {
		log.Info("Docker image: ", *image, " does NOT exist in registry: ", *registry)
	}
}
