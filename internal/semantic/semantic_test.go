package semantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagWithoutChars(t *testing.T) {
	act, err := tagWithoutChars("1.2.3-helloworld")
	if err != nil {
		t.Error(err)
	}
	exp := 1002003
	assert.Equal(t, exp, act)
}

func TestTagWithoutCharsError(t *testing.T) {
	act, err := tagWithoutChars("helloworld")
	exp := "no match was found. Verify whether the semanticTag: 'helloworld' contains characters"
	assert.EqualError(t, err, exp)
	assert.Equal(t, 0, act)
}

func TestToInt(t *testing.T) {
	act, err := toInt("4.5.6")
	if err != nil {
		t.Error(err)
	}
	exp := 4005006
	assert.Equal(t, exp, act)
}

func TestToIntError(t *testing.T) {
	act, err := toInt("1.2.3-helloworld")
	exp := "strconv.Atoi: parsing \"3-helloworld\": invalid syntax"
	assert.EqualError(t, err, exp)
	assert.Equal(t, 0, act)
}

func TestSortAndGetLatestTag(t *testing.T) {
	m := make(map[int]string)
	m[123] = "hello"
	m[789] = "boo"
	m[456] = "world"
	act := sortAndGetLatestTag(m)
	exp := "boo"
	assert.Equal(t, exp, act)
}

func TestUpdateJava(t *testing.T) {
	exp := 9
	act, err := update("16.0.1_9-jre-hotspot-focal")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, exp, act)
}

func TestUpdateGolangAlpine(t *testing.T) {
	exp := 14
	act, err := update("1.17.0-alpine3.14")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, exp, act)
}
