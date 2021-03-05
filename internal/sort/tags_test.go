package sort

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemantic(t *testing.T) {
	act, err := semantic("1.2.3")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, act)
}

func TestTags(t *testing.T) {
	exp := "11.0.1"
	act, err := Tags([]string{"1.2.3", "1.2.300", "latest", "10.9.0", "11.0.1", "10.0.0", "10.0.0-alpine", "7.80.9", "7.8.9", "4.5.6"})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, exp, act)
}
