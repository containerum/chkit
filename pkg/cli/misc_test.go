package cli

import (
	"testing"

	"github.com/blang/semver"
)

func init() {
	if API_ADDR == "" {
		panic("[INIT] API_ADDR not defined")
	}
}

func TestSemverString(test *testing.T) {
	_, err := semver.Parse(Version)
	if err != nil {
		test.Fatalf("invalid semver string: %v", err)
	}
}
