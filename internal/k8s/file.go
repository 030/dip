package k8s

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var files = make(map[string][]byte)

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		ext := filepath.Ext(s)
		if (ext == ".yml") || (ext == ".yaml") {
			b, err := os.ReadFile(filepath.Clean(s))
			if err != nil {
				return err
			}
			files[s] = b
		}
	}
	return nil
}

func FileTag(image string) (string, error) {
	if err := filepath.WalkDir(".", walk); err != nil {
		return "", err
	}

	tag := ""
	for file, content := range files {
		r, err := regexp.Compile(`image: ` + image + `:([a-z0-9\.]+)`)
		if err != nil {
			return "", err
		}
		if !r.Match(content) {
			log.Debugf("Image: '%s' not found in k8sfile: '%s'", image, file)
		} else {
			tag = string(r.FindSubmatch(content)[1])
			log.Infof("Image: '%s' tag: '%s' found in k8sfile: '%s'", image, tag, file)
		}
	}

	return tag, nil
}
