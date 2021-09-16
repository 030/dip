package gitactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"

	log "github.com/sirupsen/logrus"
)

type Elements struct {
	GitFile, GitProjectID, GitURL, RegexToFindLatestImageTag, Pass, Reviewer, Tag, TargetBranch, User string
}

type Payload struct {
	Title        string `json:"title"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	AssigneeID   string `json:"assignee_id"`
}

func (e *Elements) clone(gitCheckoutRepositoryPath, gitURL string) error {
	log.Infof("Git cloning: '%s'...", gitURL)
	_, err := git.PlainClone(gitCheckoutRepositoryPath, false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: e.User,
			Password: e.Pass,
		},
		Depth: 1,
		URL:   gitURL,
	})
	if err != nil {
		return err
	}
	return nil
}

func gitWorkTree(sourceBranch, gitCheckoutRepositoryPath string) (*git.Worktree, error) {
	r, err := git.PlainOpen(gitCheckoutRepositoryPath)
	if err != nil {
		return nil, err
	}

	//

	log.Info("git branch", sourceBranch)

	headRef, err := r.Head()
	if err != nil {
		return nil, err
	}

	branch := plumbing.ReferenceName("refs/heads/" + sourceBranch)

	// Create a new plumbing.HashReference object with the name of the branch
	// and the hash from the HEAD. The reference name should be a full reference
	// name and not an abbreviated one, as is used on the git cli.
	//
	// For tags we should use `refs/tags/%s` instead of `refs/heads/%s` used
	// for branches.
	ref := plumbing.NewHashReference(branch, headRef.Hash())

	// The created reference is saved in the storage.
	err = r.Storer.SetReference(ref)
	if err != nil {
		return nil, err
	}

	//

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	//

	w.Checkout(&git.CheckoutOptions{
		Branch: branch,
	})

	return w, nil
}

func gitNumberOfChangedFiles(sourceBranch, gitCheckoutRepositoryPath string) (int, error) {
	log.Info("Checking git status...")

	w, err := gitWorkTree(sourceBranch, gitCheckoutRepositoryPath)
	if err != nil {
		return 0, err
	}

	status, err := w.Status()
	if err != nil {
		return 0, err
	}
	return len(status), nil
}

func gitAddAndCommit(gitCheckoutRepositoryPath string) error {
	log.Info("Adding and committing changed files...")

	r, err := git.PlainOpen(gitCheckoutRepositoryPath)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	commit, err := w.Commit("AutoUpdate various versions", &git.CommitOptions{
		Author: &object.Signature{
			Email: "dip@030.github.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	o, err := r.CommitObject(commit)
	if err != nil {
		return err
	}
	log.Info(o)

	return nil
}

func (e *Elements) push(gitCheckoutRepositoryPath string) error {
	log.Info("Pushing to git...")
	r, err := git.PlainOpen(gitCheckoutRepositoryPath)
	if err != nil {
		return err
	}
	err = r.Push(&git.PushOptions{Auth: &githttp.BasicAuth{
		Username: e.User,
		Password: e.Pass,
	}})
	if err != nil {
		return err
	}
	return nil
}

//
//
//
//
//
func updateImage(file, imageRegex, latestTag string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Warningf("no kubernetes file found, as file: '%s' does not exist", file)
		return nil // continue
	}

	log.Infof("Updating image tag in k8s file: '%s'...", file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	matched, err := regexp.Match(imageRegex, b)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("unable to update docker tag. Verify whether the regex: '%s' matches the FROM statement in the Dockerfile: '%s'", file, string(b))
	}

	re, err := regexp.Compile(imageRegex)
	if err != nil {
		return err
	}
	replaced := re.ReplaceAll(b, []byte("${1}"+latestTag+"${2}"))
	if err = ioutil.WriteFile(file, replaced, 0644); err != nil {
		return err
	}

	return nil
}

func (e *Elements) CloneAndUpdateTags(image, latestTag, sourceBranch string) error {
	log.Infof("Repository: '%s'", e.GitURL)
	t := time.Now().UnixNano()
	s := strconv.FormatInt(t, 10)
	gitCheckoutRepositoryPath := filepath.Join("/tmp", "dip", "repositories", s)

	if err := e.clone(gitCheckoutRepositoryPath, e.GitURL); err != nil {
		return err
	}

	// if ce.DockerFileRegexReplace != "" {
	// 	dockerfilePath := filepath.Join(gitCheckoutRepositoryPath, "Dockerfile")
	// 	if err := updateFROMStatementDockerfile(dockerfilePath, ce.DockerFileRegexReplace, ce.LatestTag); err != nil {
	// 		return err
	// 	}
	// }

	// if ce.Name == "golang" {
	// 	goModfilePath := filepath.Join(gitCheckoutRepositoryPath, "go.mod")
	// 	if err := updateGoModVersion(goModfilePath, ce.LatestTag); err != nil {
	// 		return err
	// 	}

	// 	githubRepo, err := regexp.MatchString(`github\.com`, repository)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if githubRepo {
	// 		if _, err := os.Stat(filepath.Join(gitCheckoutRepositoryPath, ".github", "workflows")); !os.IsNotExist(err) {
	// 			if err := updateGoVersionInActions(gitCheckoutRepositoryPath, ce.LatestTag); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

	if err := updateImage(filepath.Join(gitCheckoutRepositoryPath, e.GitFile), e.RegexToFindLatestImageTag, latestTag); err != nil {
		return err
	}

	numberOfChanges, err := gitNumberOfChangedFiles(sourceBranch, gitCheckoutRepositoryPath)
	if err != nil {
		return err
	}
	if numberOfChanges > 0 {
		if err := gitAddAndCommit(gitCheckoutRepositoryPath); err != nil {
			return err
		}
		if err := e.push(gitCheckoutRepositoryPath); err != nil {
			return err
		}
	} else {
		log.Infof("Skipping git add and commit for: '%s', in: '%s' as nothing was changed. Number of changes: '%d'", "bla", e.GitURL, numberOfChanges)
	}

	//
	//
	// MR
	//
	//

	data := Payload{
		Title:        sourceBranch,
		SourceBranch: sourceBranch,
		TargetBranch: e.TargetBranch,
		AssigneeID:   e.Reviewer,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://gitlab.com/api/v4/projects/"+e.GitProjectID+"/merge_requests", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(e.User, e.Pass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Info(resp.StatusCode)

	return nil
}

//
//
//
// todo
//
//
//

func (e *Elements) CreateMR(image, latestTag string) error {
	sourceBranch := image + "-" + latestTag
	if err := e.CloneAndUpdateTags(image, latestTag, sourceBranch); err != nil {
		return err
	}
	return nil
}
