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
	got := tags("mongo").StatusCode
	want := 200
	if got != want {
		t.Errorf("Want '%v', got %v", want, got)
	}
}
