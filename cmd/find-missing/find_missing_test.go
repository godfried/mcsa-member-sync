package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestFindMissing(t *testing.T) {
	resultsDir, err := ioutil.TempDir("", "mcsa-member-sync-find-missing-tests-*")
	assert.NoError(t, err)
	for _, test := range []string{"test1"} {
		fixtureDir := filepath.Join("fixtures", test)
		sm := filepath.Join(fixtureDir, "source_membaz.csv")
		se := filepath.Join(fixtureDir, "source_everlytic.csv")
		dm := filepath.Join(resultsDir, "destination_membaz.csv")
		de := filepath.Join(resultsDir, "destination_everlytic.csv")
		err := run(sm, se, dm, de)
		assert.NoError(t, err)
		em := filepath.Join(fixtureDir, "expected_membaz.csv")
		ee := filepath.Join(fixtureDir, "expected_everlytic.csv")
		for expectedFile, actualFile := range map[string]string{
			em: dm,
			ee: de,
		} {
			expected, err := ioutil.ReadFile(expectedFile)
			assert.NoError(t, err)
			actual, err := ioutil.ReadFile(actualFile)
			assert.NoError(t, err)
			assert.Equal(t, string(expected), string(actual))
		}
	}
}
