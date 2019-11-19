package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/buger/jsonparser"
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

// tags returns tag json formatted information from dockerhub
func tags(image string) *http.Response {
	// if image does not contain a forward slash, the assumption is that it
	// is a library
	if !strings.Contains(image, "/") {
		image = "library/" + image
	}

	resp, err := http.Get("https://registry.hub.docker.com/v2/repositories/" + image + "/tags?page_size=1024")
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

// readResp reads a http response and returns it as byte
func readResp(resp *http.Response) []byte {
	defer resp.Body.Close()
	var b []byte
	var err error
	if resp.StatusCode == http.StatusOK {
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	return b
}

// latestTag returns the latest tag of a docker image
func latestTag(b []byte, t string) string {
	var c string
	var arr []string
	jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		s, _ := jsonparser.GetString(value, "name")
		a := fmt.Sprintf(`%s`, t)
		re := regexp.MustCompile(a)
		if re.FindString(s) != "" {
			c = fmt.Sprintf("%v", re.FindString(s))
			arr = append(arr, c)
		}
	}, "results")
	return arr[0]
}

func main() {
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
	image := flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
	latest := flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")
	registry := flag.String("registry", "", "To what destination the image should be transferred, e.g. quay.io/some-org/. Note: do not omit the last forward slash.")
	preserve := flag.Bool("preserve", false, "Whether an image from dockerhub should be stored in a private registry.")
	date := flag.Bool("date", false, "Sometimes the version of an image gets overwritten by the community due to security updates. In order to store the latest image in a private registry, one could append a date.")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	var latestDetectedTag string
	if *latest != "" {
		latestDetectedTag = ":" + latestTag(readResp(tags(*image)), *latest)
		fmt.Println(latestDetectedTag)
	}

	i := *image + latestDetectedTag
	// forward slashes are not allowed in some registries like quay.io,
	// e.g. sonatype/nexus3:3.19.1 will become sonatype-nexus3:3.19.1
	i = strings.Replace(i, "/", "-", -1)

	var d string
	if *date {
		currentTime := time.Now()
		d = i + "-" + currentTime.Format("20060102-150405")
	} else {
		d = i
	}

	if *registry != "" {
		dockerImageAbsent := absent(d, *registry)
		if !dockerImageAbsent {
			msg := "Docker image: " + d + " already exists in registry: " + *registry
			if *preserve {
				// Never return an exit1 if the aim is to preserve an image as
				// the CI will become RED, while it should be green if an image
				// is already present
				log.Info(msg)
				os.Exit(0)
			} else {
				// Return an Exit1 if an image already exists to prevent that
				// it gets overwritten if tagImmutability is absent in a
				// docker registry
				log.Fatal(msg)
			}
		} else {
			log.Info("Docker image: ", d, " does NOT exist in registry: ", *registry)
		}
	}

	if *preserve {
		var cmd string

		cmd = "docker pull " + i
		log.Info(cmd)
		command(cmd)

		cmd = "docker tag " + i + " " + *registry + d
		log.Info(cmd)
		err := command(cmd)
		if err != nil {
			log.Fatal(err)
		}

		cmd = "docker push " + *registry + d
		log.Info(cmd)
		err = command(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
}
