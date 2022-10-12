package quay

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatest(t *testing.T) {
	json, err := os.ReadFile(filepath.Join("..", "..", "test", "testdata", "quay.json"))
	if err != nil {
		t.Error(err)
	}

	exp := "v1.9.1"
	act, err := latest(json, `^v1(\.[0-9]+){2}$`)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exp, act)
}

func TestLatestFail(t *testing.T) {
	json, err := os.ReadFile(filepath.Join("..", "..", "test", "testdata", "quay.json"))
	if err != nil {
		t.Error(err)
	}

	_, err = latest(json, `^v3(\.[0-9]+){2}$`)
	assert.EqualError(t, err, "no tags were found. Check whether regex is correct")
}
