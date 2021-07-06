package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
			err := unite([]string{src, dest})
			if err != nil {
				t.Fatalf("Unexpected error uniting template: %v", err)
			}

			testUniteOutput(t, tc.Name(), dest)
		})
	}
}

func TestUniteTemplate(t *testing.T) {
	testCases, err := os.ReadDir("testdata/templated")
	if err != nil {
		t.Fatalf("Unable to read testdata/templated dir: %v", err)
	}

	for _, tc := range testCases {
		// Check for each data file we find
		dir := filepath.Join("testdata/templated", tc.Name())
		fs, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("Unable to read directory %q: %v", dir, err)
		}
		for _, f := range fs {
			ext := filepath.Ext(f.Name())
			if ext != ".json" && ext != ".yaml" && ext != ".yml" {
				continue
			}

			dataFile := filepath.Join(dir, f.Name())
			t.Run(filepath.Join(tc.Name(), f.Name()), func(t *testing.T) {
				src := filepath.Join("testdata/templated", tc.Name())
				dest := filepath.Join(t.TempDir(), "template.yml")
				err := unite([]string{src, dest, "--data-file", dataFile})
				if err != nil {
					t.Fatalf("Unexpected error uniting template: %v", err)
				}

				testUniteOutput(t, tc.Name(), dest)
			})
		}
	}
}

func testUniteOutput(t *testing.T, dir, dest string) {
	exp := loadOutput(t, filepath.Join("testdata/united", dir, "template.yml"))
	act := loadOutput(t, dest)
	if diff := cmp.Diff(exp, act, cmp.Comparer(cmpStrings)); diff != "" {
		t.Errorf("File contents different:\n%s", diff)
	}
}

func loadOutput(t *testing.T, filename string) []interface{} {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("unexpected error opening %q: %v", filename, err)
	}

	// We have
	var buf bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
		buf.WriteString("\n")
	}

	var items []interface{}
	dec := yaml.NewDecoder(&buf)
	for {
		var item interface{}
		if err := dec.Decode(&item); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("unable to decode item in %q: %v", filename, err)
		}

		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return getItemName(items[i]) < getItemName(items[j])
	})
	return items
}

func getItemName(item interface{}) string {
	return item.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
}

// cmpStrings is needed to compare the strings in the yaml, ignoring any trailing whitespace/line ends.
func cmpStrings(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}
