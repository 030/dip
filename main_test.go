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
