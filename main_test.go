package main

import (
	"testing"
)

func TestCommand(t *testing.T) {
	// happy
	got := command("ls")
	if got != nil {
		t.Errorf("Want 'no error', got %v", got)
	}

	// unhappy
	got = command("hello")
	want := "Cannot run: 'hello'. Error: 'bash: hello: command not found\n', exit status 127"
	if want != got.Error() {
		t.Errorf("Want %v, got %v", want, got)
	}
}

func TestTags(t *testing.T) {
	// happy
	got, _ := tags("mongo")
	gotInt := len(got)
	wantGreaterThan := 0
	if !(gotInt > wantGreaterThan) {
		t.Errorf("Mongo exists in dockerhub and thus the number of tags should be greater than 0. Want greater than '%d', got '%d'.", wantGreaterThan, gotInt)
	}

	// unhappy
	_, err := tags("image-does-not-exist")
	want := "No versions were found. Check whether image 'library/image-does-not-exist' exists in the registry"
	if want != err.Error() {
		t.Errorf("Want '%v', got '%v'", want, err)
	}
}

func TestLatestTags(t *testing.T) {
	tagsSemantic := []string{"3.0.0", "1.1.0", "1.1.42", "2.0.0", "1.0.0"}
	tagsNonSemantic := []string{"19.04", "19.100", "20.04", "42.08", "19.10"}

	// happy - semantic
	got, _ := latestTag(tagsSemantic, "^1.1.*", true)
	want := "1.1.42"
	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}

	// happy - non-semantic
	got, _ = latestTag(tagsNonSemantic, "^19.*", false)
	want = "19.100"
	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}

	// unhappy - semantic
	_, err := latestTag(tagsSemantic, "^9.*", true)
	want = "None of the tags: [3.0.0 1.1.0 1.1.42 2.0.0 1.0.0] match regex: ^9.*"
	if want != err.Error() {
		t.Errorf("Want '%v', got '%v'", want, err)
	}

	// unhappy - non-semantic
	_, err = latestTag(tagsNonSemantic, "^9.*", false)
	want = "None of the tags: [19.04 19.100 20.04 42.08 19.10] match regex: ^9.*"
	if want != err.Error() {
		t.Errorf("Want '%v', got '%v'", want, err)
	}
}
