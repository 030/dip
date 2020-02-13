package main

import "testing"

func TestCommand(t *testing.T) {
	// happy
	got := command("lls")
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
