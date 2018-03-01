package cmd

import (
	"testing"

	"github.com/blang/semver"
)

func TestSemverString(test *testing.T) {
	_, err := semver.Parse(Version)
	if err != nil {
		test.Fatalf("invalid semver string: %v", err)
	}
}
