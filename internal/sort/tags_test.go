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

func TestSomething(t *testing.T) {
	// [latest stable-alpine-perl stable-alpine 1.18.0-alpine-perl 1.18.0-alpine 1.18-alpine-perl 1.18-alpine mainline-alpine-perl mainline-alpine alpine-perl alpine 1.19.7-alpine-perl 1.19.7-alpine 1.19-alpine-perl 1.19-alpine 1-alpine-perl 1-alpine 1 perl mainline-perl mainline 1.19.7-perl 1.19.7 1.19-perl 1.19 1-perl stable-perl stable 1.19.6-perl 1.19.6 1.18.0-perl 1.18.0 1.18-perl 1.18 1.19.6-alpine-perl 1.19.6-alpine 1.19.5-perl 1.19.5-alpine-perl 1.19.5-alpine 1.19.5 1.19.4-perl 1.19.4 1.19.4-alpine-perl 1.19.4-alpine 1.19.3-alpine-perl 1.19.3-alpine 1.19.3-perl 1.19.3 1.19.2-perl 1.19.2 1.19.2-alpine-perl 1.19.2-alpine 1.19.1-perl 1.19.1 1.19.1-alpine-perl 1.19.1-alpine 1.19.0-perl 1.19.0 1.19.0-alpine-perl 1.19.0-alpine 1.17.10-perl 1.17.10 1.17-perl 1.17 1.17.10-alpine-perl 1.17.10-alpine 1.17-alpine-perl 1.17-alpine 1.16.1-perl 1.16.1 1.16-perl 1.16 1.17.9-perl 1.17.9 1.17.9-alpine-perl 1.17.9-alpine 1.17.8-perl 1.17.8 1.17.8-alpine-perl 1.17.8-alpine 1.16.1-alpine-perl 1.16.1-alpine 1.16-alpine-perl 1.16-alpine 1.17.7-perl 1.17.7-alpine-perl 1.17.7-alpine 1.17.7 1.17.6-perl 1.17.6 1.17.6-alpine-perl 1.17.6-alpine 1.17.5-alpine-perl 1.17.5-alpine 1.17.5-perl 1.17.5 1.17.4-alpine-perl 1.17.4-alpine 1.17.4-perl 1.17.4 1.17.3-perl 1.17.3 1.17.3-alpine-perl 1.17.3-alpine 1.17.2-perl 1.17.2 1.16.0-perl 1.16.0 1.17.2-alpine-perl 1.17.2-alpine 1.17.1-perl 1.17.1 1.17.1-alpine-perl 1.17.1-alpine 1.17.0-alpine-perl 1.17.0-alpine 1.16.0-alpine-perl 1.16.0-alpine 1.17.0-perl 1.17.0 1.15.12-alpine-perl 1.15.12-alpine 1.15-alpine-perl 1.15-alpine 1.15.12-perl 1.15.12 1.15-perl 1.15 1.15.11-alpine-perl 1.15.11-alpine 1.15.11-perl 1.15.11 1.14-alpine-perl 1.14.2-alpine-perl 1.14-alpine 1.14.2-alpine 1.15.10-alpine-perl 1.15.10-alpine 1.15.10-perl 1.15.10 1.14.2 1.14-perl 1.14.2-perl 1.14 1.15.9-alpine-perl 1.15.9-alpine 1.15.9-perl 1.15.9 1.15.8-perl 1.15.8 1.15.8-alpine-perl 1.15.8-alpine 1.15.7-alpine-perl 1.15.7-alpine 1.15.7-perl 1.15.7 1.14.1-alpine-perl 1.14.1-alpine 1.15.6-alpine-perl 1.15.6-alpine 1.14.1-perl 1.14.1 1.15.6-perl 1.15.6 1.15.5-perl 1.15.5-alpine-perl 1.14.0-perl 1.14.0 1.15.5 1.14.0-alpine-perl 1.14.0-alpine 1.15.5-alpine 1.15.4-alpine-perl 1.15.4-alpine 1.15.4-perl 1.15.4 1.15.3-alpine-perl 1.15.3-alpine 1.15.3-perl 1.15.3 1.15.2-alpine-perl 1.15.2 1.15.2-alpine 1.15.2-perl 1.15.1-alpine 1.15.1-perl 1.15.1 1.15.1-alpine-perl 1.15.0-alpine 1.15.0-perl 1.15.0 1.15.0-alpine-perl 1.13-alpine-perl 1.13.12-alpine-perl 1.13-alpine 1.13.12-alpine 1.13-perl 1.13.12-perl 1.13 1.13.12 1.12-perl 1.12.2-perl 1.12 1.12.2 1.13.11-alpine-perl 1.13.11-alpine 1.13.11-perl 1.13.11 1.13.10-alpine-perl 1.13.10-alpine 1.13.10-perl 1.13.10 1.12-alpine 1.13.9-perl 1.13.9 1.13.9-alpine 1.13.9-alpine-perl 1.13.8-perl 1.13.8 1.12-alpine-perl 1.12.2-alpine-perl 1.12.2-alpine 1.13.8-alpine-perl 1.13.8-alpine 1.13.7-perl 1.13.7 1.13.7-alpine-perl 1.13.7-alpine 1.13.6-perl 1.13.6 1.13.6-alpine-perl 1.13.6-alpine 1.12.1-alpine-perl 1.12.1-alpine 1.13.5-alpine-perl 1.13.5-alpine 1.12.1-perl 1.12.1 1.13.5-perl 1.13.5 1.13.3-alpine-perl 1.13.3-alpine 1.13.3-perl 1.13.3 1.12.0-alpine-perl 1.12.0-alpine 1.12.0-perl 1.12.0 1.13.2-alpine-perl 1.13.2-alpine 1.13.2-perl 1.13.2 1.13.1-alpine-perl 1.13.1-alpine 1.13.1-perl 1.13.1 1.13.0-perl 1.13.0-alpine-perl 1.13.0-alpine 1.13.0 1.11-alpine 1.11.13-alpine 1.11 1.11.13 1.10-alpine 1.10.3-alpine 1.10 1.10.3 1.11.12-alpine 1.11.12 1.11.10-alpine 1.11.10 1.11.9-alpine 1.11.9 1.10.2-alpine 1.10.2 1.11.8 1.11.8-alpine 1.11.7-alpine 1.11.7 1.11.6-alpine 1.11.6 1.11.5 1.11.5-alpine 1.10.1-alpine 1.10.1 1.11.4 1.11.4-alpine 1.11.3-alpine 1.11.3 1.11.1-alpine 1.11.0-alpine 1.11.0 1.11.1 1.10.0-alpine 1.10.0 1.9.15-alpine 1.9-alpine 1.9.15 1.9 1.8.1-alpine 1.8-alpine 1.9.14-alpine 1.8.1 1.8 1.9.14 1.9.12 1.9.11 1.9.10 1.9.9 1.9.8 1.9.7 1.7.8 1.7.7 1.7.6 1.7.9 1.7.10 1.7.1 1.7.5 1.9.6 1.9.5 1.9.4 1.9.3 1.9.2 1.9.0 1.9.1 1.7.11 1.7.12 1.7]
	act, err := something([]string{"1.2.3", "1.2.300", "latest", "10.9.0", "10.0.0", "10.0.0-alpine", "7.80.9", "7.8.9", "4.5.6"})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "50.51.600", act)
}
