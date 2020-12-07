package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/go-version"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// https://stackoverflow.com/a/11355611/2777965
var ver string

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
		log.Debug("Image: '" + image + "' is a library. Concatenating 'library/'...")
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

// latestTag returns the latest tag of a docker image
func latestTag(s []string, t string, z bool) (string, error) {
	log.Debug("Getting all tags that match regex: '", t, "'")
	var c string
	var arr []string

	re := regexp.MustCompile(t)
	for _, x := range s {
		if re.FindString(x) != "" {
			c = fmt.Sprintf("%v", re.FindString(x))
			log.Debug(c)
			arr = append(arr, c)
			log.Debug(arr)
		}
	}

	if len(arr) == 0 {
		return "", fmt.Errorf("None of the tags: %v match regex: %s", s, t)
	}

	var versions []*version.Version
	var latestVersionString string
	if z {
		log.Debug("Raw slice latestTag: ", arr)
		// Following snippet retrieved from
		// https://github.com/hashicorp/go-version#version-sorting
		versions = make([]*version.Version, len(arr))
		for i, raw := range arr {
			v, _ := version.NewVersion(raw)
			versions[i] = v
		}

		sort.Sort(version.Collection(versions))
		log.Debug("Sorted slice latestTag: ", versions)
		// Retrieve last element, see https://stackoverflow.com/a/22535888
		latestVersion := versions[len(versions)-1]
		latestVersionString = latestVersion.String()
	} else {
		sort.Strings(arr)
		log.Debug("Sorted slice:", arr)
		latestVersionString = arr[len(arr)-1]
	}

	return latestVersionString, nil
}

func main2() {
	version := flag.Bool("version", false, "Return the version of the tool.")
	debug := flag.Bool("debug", false, "Whether debug mode should be enabled.")
	semantic := flag.Bool("semantic", true, "Whether the tags are semantic.")
	image := flag.String("image", "", "Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.")
	latest := flag.String("latest", "", "The regex to get the latest tag, e.g. \"xenial-\\d.*\".")
	registry := flag.String("registry", "", "To what destination the image should be transferred, e.g. quay.io/some-org/. Note: do not omit the last forward slash.")
	preserve := flag.Bool("preserve", false, "Whether an image from dockerhub should be stored in a private registry.")
	date := flag.Bool("date", false, "Sometimes the version of an image gets overwritten by the community due to security updates. In order to store the latest image in a private registry, one could append a date.")

	flag.Parse()

	if *version {
		fmt.Println("dip version " + ver)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithFields(log.Fields{
		"debug":    *debug,
		"semantic": *semantic,
		"image":    *image,
		"latest":   *latest,
		"registry": *registry,
		"preserve": *preserve,
		"date":     *date,
	}).Debug("Docker Image Patrol (DIP) command line arguments:")

	var latestDetectedTag string
	if *latest != "" {

		t, err := tags(*image)
		if err != nil {
			log.Fatal(err)
		}

		l, err := latestTag(t, *latest, *semantic)
		if err != nil {
			log.Fatal(err)
		}

		latestDetectedTag = ":" + l
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
				log.Debug(msg)
				os.Exit(0)
			} else {
				// Return an Exit1 if an image already exists to prevent that
				// it gets overwritten if tagImmutability is absent in a
				// docker registry
				log.Fatal(msg)
			}
		} else {
			log.Warning("Docker image: ", d, " does NOT exist in registry: ", *registry)
		}
	}

	if *preserve {
		var cmd string

		cmd = "docker pull " + *image + latestDetectedTag
		log.Info(cmd)
		err := command(cmd)
		if err != nil {
			log.Fatal(err)
		}

		cmd = "docker tag " + *image + latestDetectedTag + " " + *registry + d
		log.Info(cmd)
		err = command(cmd)
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

var hihi = []string{}

func boo(s string, page int) error {
	resp, err := grequests.Get("https://registry.hub.docker.com/v2/repositories/"+s+"/tags?page="+strconv.Itoa(page)+"&page_size=100", nil)
	if err != nil {
		return err
	}
	httpStatusCode := resp.StatusCode
	if httpStatusCode != http.StatusOK {
		return fmt.Errorf("ResponseCode not 200, but: '%v'. Check whether image: '%v', exists on dockerhub", httpStatusCode, s)
	}

	hihi = append(hihi, bladibla(resp.Bytes())...)
	if gjson.GetBytes(resp.Bytes(), "next").String() != "" {
		fmt.Println(gjson.GetBytes(resp.Bytes(), "next"))
		page++
		fmt.Println(page)
		if err := boo(s, page); err != nil {
			return err
		}
	}
	return nil
}

func bladibla(b []byte) []string {
	tags := gjson.GetBytes(b, "results.#.name").Array()
	boo := []string{}
	for _, tag := range tags {
		boo = append(boo, tag.String())
	}
	return boo
}

func main() {
	if err := boo("library/tomcat", 1); err != nil {
		log.Fatal(err)
	}
	sort.Sort(sort.StringSlice(hihi))
	fmt.Println(hihi)
	fmt.Println(len(hihi))
}
