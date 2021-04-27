package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestSplit(t *testing.T) {
	testCases, err := os.ReadDir("testdata/united")
	if err != nil {
		t.Fatalf("Unable to read testdata/split dir: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			dir := t.TempDir()
			filename := filepath.Join("testdata/united", tc.Name(), "template.yml")
			err := split([]string{filename, dir})
			if err != nil {
				t.Fatalf("Unexpected error splitting template: %v", err)
			}

			compareDirs(t, dir, filepath.Join("testdata/split", tc.Name()))

		})
	}
}

// compareDirs compares the contents of two directories and makes sure
// that they are both the same
func compareDirs(t *testing.T, dir0, dir1 string) {
	t.Helper()

	ch0, ch1 := make(chan string), make(chan string)
	go walkDir(t, dir0, ch0)
	go walkDir(t, dir1, ch1)

	for {
		filename0, ok0 := <-ch0
		filename1, ok1 := <-ch1

		if !ok0 || !ok1 {
			if ok0 != ok1 {
				t.Errorf("different number of files found")
			}
			break
		}

		if filepath.Base(filename0) != filepath.Base(filename1) {
			t.Fatalf("different files found: %q and %q", filename0, filename1)
		}

		b0, err := os.ReadFile(filename0)
		if err != nil {
			t.Fatalf("unexpected error opening %q: %s", filename0, err)
		}
		b1, err := os.ReadFile(filename1)
		if err != nil {
			t.Fatalf("unexpected error opening %q: %s", filename0, err)
		}

		// If we have two templates, compare the yaml
		if filepath.Ext(filename0) == ".yml" {
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
			return
		}

		// Otherwise, compare the strings
		s0 := strings.TrimSpace(string(b0))
		s1 := strings.TrimSpace(string(b1))
		if diff := cmp.Diff(s0, s1); diff != "" {
			t.Errorf("File contents different:\n%s", diff)
		}
	}
}

func walkDir(t *testing.T, dir string, ch chan string) {
	defer close(ch)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ch <- path
		return nil
	})
	if err != nil {
		t.Errorf("unable to walk dir %q: %v", dir, err)
	}
}
