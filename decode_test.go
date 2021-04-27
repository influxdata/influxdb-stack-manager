package main

import (
	"os"
	"testing"
)

func TestDecode(t *testing.T) {

	testCases := []string{
		"testdata/dashboard.yaml",
		"testdata/task.yaml",
		"testdata/label.yaml",
		"testdata/multiple_templates.yaml",
	}

	os.Remove("testoutput")
	os.Mkdir("testoutput", 0700)

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			f, err := os.Open(tc)
			if err != nil {
				t.Fatalf("error opening file: %v", err)
			}
			defer f.Close()

			err = decode("testoutput", f)
			t.Error(err)
		})
	}
}
