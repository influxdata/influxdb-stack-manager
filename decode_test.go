package main

import (
	"os"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {

	testCases := []string{
		// "testdata/dashboard.yaml",
		// "testdata/task.yaml",
		// "testdata/label.yaml",
		"testdata/multiple_templates.yaml",
	}

	os.Remove("testoutput")
	os.Remove("testtmp")
	os.Mkdir("testoutput", 0700)
	os.Mkdir("testtmp", 0700)

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			f, err := os.Open(tc)
			if err != nil {
				t.Fatalf("error opening file: %v", err)
			}
			defer f.Close()

			// dir := t.TempDir()
			if err := split("testtmp", f); err != nil {
				t.Fatalf("Unable to split template: %v", err)
			}

			outputFile := strings.Replace(tc, "testdata", "testoutput", 1)
			f, err = os.Create(outputFile)
			if err != nil {
				t.Fatalf("error creating output file: %v", err)
			}
			defer f.Close()

			if err := unite("testtmp", f); err != nil {
				t.Fatalf("unable to unite files: %v", err)
			}
		})
	}
}
