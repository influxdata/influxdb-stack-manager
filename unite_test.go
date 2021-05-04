package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestUnite(t *testing.T) {
	testCases, err := os.ReadDir("testdata/split")
	if err != nil {
		t.Fatalf("Unable to read testdata/split dir: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			src := filepath.Join("testdata/split", tc.Name())
			dest := filepath.Join(t.TempDir(), "template.yml")
			unite([]string{src, dest})

			exp := filepath.Join("testdata/united", tc.Name(), "template.yml")
			b0, err := os.ReadFile(exp)
			if err != nil {
				t.Fatalf("unexpected error opening %q: %s", exp, err)
			}
			b1, err := os.ReadFile(dest)
			if err != nil {
				t.Fatalf("unexpected error opening %q: %s", dest, err)
			}

			var tmpl0, tmpl1 interface{}
			if err := yaml.Unmarshal(b0, &tmpl0); err != nil {
				t.Fatalf("unexpected error unmarshing b0: %v", err)
			}
			if err := yaml.Unmarshal(b1, &tmpl1); err != nil {
				t.Fatalf("unexpected error unmarshing b1: %v", err)
			}

			if diff := cmp.Diff(tmpl0, tmpl1); diff != "" {
				t.Errorf("File contents different:\n%s", diff)
			}
		})
	}
}
