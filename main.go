package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
)

func command(s string) error {
	out, err := exec.Command("bash", "-c", s).CombinedOutput()
	outString := string(out)

	log.Debug(s)
	log.Info(outString)

	if err != nil {
		return fmt.Errorf("Cannot run: '%s'. Error: '%s', %v", s, outString, err)
	}
	return nil
}

// tags returns tag json formatted information from dockerhub
func tags(image string) ([]string, error) {
	// if image does not contain a forward slash, the assumption is that it
	// is a library
	log.Debug("Checking whether image: '" + image + "' is a library")
	if !strings.Contains(image, "/") {
		log.Info("Image: '" + image + "' is a library. Concatenating 'library/'...")
		image = "library/" + image
	}

	log.Debug("Getting raw tag information on dockerhub for image: '" + image + "'")

	// get url
	resp, err := http.Get("https://registry.hub.docker.com/v2/repositories/" + image + "/tags?page_size=1024")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body into bytes
	var b []byte
	if resp.StatusCode == http.StatusOK {
		log.Debug("Dockerhub api call Ok. Reading body as []byte")
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	// get tags and return them as a slice
	// var c string
	var arr []string
	_, err = jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		key := "name"

		s, _ := jsonparser.GetString(value, key)
		if s == "" {
			log.Warning("No value retrieved for key: '" + key + "'")
		}

		arr = append(arr, s)
		log.Debug(arr)

	}, "results")

	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("No versions were found. Check whether image '" + image + "' exists in the registry")
	}
	return arr, nil
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

func sortedLatest(s []string) string {
	versionsRaw := s

	// Following snippet retrieved from
	// https://github.com/hashicorp/go-version#version-sorting
	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))

	// Retrieve last element, see https://stackoverflow.com/a/22535888
	latestVersion := versions[len(versions)-1]
	return latestVersion.String()
}

// latestTag returns the latest tag of a docker image
func latestTag(b []byte, t string) (string, error) {
	var c string
	var arr []string

	_, err := jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		key := "name"

		s, _ := jsonparser.GetString(value, key)
		if s == "" {
			fmt.Printf("No value retrieved for key: '%s'", key)
		}

		re := regexp.MustCompile(t)

		if re.FindString(s) != "" {
			c = fmt.Sprintf("%v", re.FindString(s))
			arr = append(arr, c)
		}
	}, "results")

	if err != nil {
		return "", err
	}

	if len(arr) == 0 {
		return "", fmt.Errorf("No versions were found. Check whether image exists in the registry")
	}

	return sortedLatest(arr), nil
}

// func main() {
// 	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
// 	image := flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
// 	latest := flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")
// 	registry := flag.String("registry", "", "To what destination the image should be transferred, e.g. quay.io/some-org/. Note: do not omit the last forward slash.")
// 	preserve := flag.Bool("preserve", false, "Whether an image from dockerhub should be stored in a private registry.")
// 	date := flag.Bool("date", false, "Sometimes the version of an image gets overwritten by the community due to security updates. In order to store the latest image in a private registry, one could append a date.")

// 	flag.Parse()

// 	if *debug {
// 		log.SetLevel(log.DebugLevel)
// 	}

// 	var latestDetectedTag string
// 	if *latest != "" {
// 		l, err := latestTag(readResp(tags(*image)), *latest)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		latestDetectedTag = ":" + l
// 		fmt.Println(latestDetectedTag)
// 	}

// 	i := *image + latestDetectedTag
// 	// forward slashes are not allowed in some registries like quay.io,
// 	// e.g. sonatype/nexus3:3.19.1 will become sonatype-nexus3:3.19.1
// 	i = strings.Replace(i, "/", "-", -1)

// 	var d string
// 	if *date {
// 		currentTime := time.Now()
// 		d = i + "-" + currentTime.Format("20060102-150405")
// 	} else {
// 		d = i
// 	}

// 	if *registry != "" {
// 		dockerImageAbsent := absent(d, *registry)
// 		if !dockerImageAbsent {
// 			msg := "Docker image: " + d + " already exists in registry: " + *registry
// 			if *preserve {
// 				// Never return an exit1 if the aim is to preserve an image as
// 				// the CI will become RED, while it should be green if an image
// 				// is already present
// 				log.Info(msg)
// 				os.Exit(0)
// 			} else {
// 				// Return an Exit1 if an image already exists to prevent that
// 				// it gets overwritten if tagImmutability is absent in a
// 				// docker registry
// 				log.Fatal(msg)
// 			}
// 		} else {
// 			log.Info("Docker image: ", d, " does NOT exist in registry: ", *registry)
// 		}
// 	}

// 	if *preserve {
// 		var cmd string

// 		cmd = "docker pull " + *image + latestDetectedTag
// 		log.Info(cmd)
// 		err := command(cmd)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		cmd = "docker tag " + *image + latestDetectedTag + " " + *registry + d
// 		log.Info(cmd)
// 		err = command(cmd)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		cmd = "docker push " + *registry + d
// 		log.Info(cmd)
// 		err = command(cmd)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }
