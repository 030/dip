package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	exp := "11.0.1"
	act, err := Tags([]string{"1.2.3", "1.2.300", "latest", "10.9.0", "11.0.1", "10.0.0", "10.0.0-alpine", "7.80.9", "7.8.9", "4.5.6"})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, exp, act)
}

func TestTagsTwo(t *testing.T) {
	exp := "5"
	act, err := Tags([]string{"5", "4", "3", "2", "1"})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, exp, act)
}

func TestTagsError(t *testing.T) {
	actualError := "cannot find the latest tag. Check whether the tags are semantic"
	_, err := Tags([]string{""})
	assert.EqualError(t, err, actualError)
}

func TestTagsTwoError(t *testing.T) {
	actualError := "tags should not be empty"
	_, err := Tags([]string{})
	assert.EqualError(t, err, actualError)
}
