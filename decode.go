package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// decode will populate the given directory with the decoded
// contents of the yaml file.
func decode(dir string, r io.Reader) error {
	seenNames := map[string]struct{}{}

	decoder := yaml.NewDecoder(r)
	for {
		var obj object
		if err := decoder.Decode(&obj); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("unable to decode template: %v", err)
		}

		var spec struct {
			Name string `yaml:"name"`
		}
		if err := obj.Spec.Decode(&spec); err != nil {
			return fmt.Errorf("unable to decode spec: %v", err)
		}

		dir := filepath.Join(dir, escapeName(spec.Name))
		// Check whether we have a name collision, just error if so.
		if _, ok := seenNames[dir]; ok {
			return fmt.Errorf("name collision detected: %q appears more than once", dir)
		}
		seenNames[dir] = struct{}{}

		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("unable to make directory %q: %v", dir, err)
		}

		var err error
		switch obj.Kind {
		case kindDashboard:
			err = writeDashboard(dir, obj)

		case kindTask:
			err = writeTask(dir, obj)

		default:
			err = writeTemplate(dir, obj)
		}
		if err != nil {
			return fmt.Errorf("unable to decode %q: %v", dir, err)
		}
	}
}

func writeDashboard(dir string, obj object) error {
	d := dashboard{
		APIVersion: obj.APIVersion,
		Kind:       obj.Kind,
		Metadata:   obj.Metadata,
	}
	if err := obj.Spec.Decode(&d.Spec); err != nil {
		return fmt.Errorf("unable to unmarshal dashboard spec: %v", err)
	}

	queryNames := map[string]int{}
	for i, chart := range d.Spec.Charts {
		for j, query := range chart.Queries {
			// Try and generate a unique name for the query
			name := escapeName(chart.Name + "_" + chart.Kind)
			filename := fmt.Sprintf("%s_%d.flux", name, queryNames[name])
			queryNames[name]++

			if err := os.WriteFile(filepath.Join(dir, filename), []byte(query.Query), 0644); err != nil {
				return fmt.Errorf("unable to write query to file %q: %v", name, err)
			}

			// Update the query in the template to just a link
			d.Spec.Charts[i].Queries[j].Query = fmt.Sprintf("file://%s", filename)
		}
	}

	return writeTemplate(dir, d)
}

func writeTask(dir string, obj object) error {
	t := task{
		APIVersion: obj.APIVersion,
		Kind:       obj.Kind,
		Metadata:   obj.Metadata,
	}
	if err := obj.Spec.Decode(&t.Spec); err != nil {
		return fmt.Errorf("unable to unmarshal task spec: %v", err)
	}

	// Write out the query to file
	if err := os.WriteFile(filepath.Join(dir, "query.flux"), []byte(t.Spec.Query), 0644); err != nil {
		return fmt.Errorf("unable to write query to file: %v", err)
	}

	t.Spec.Query = "file://query.flux"
	return writeTemplate(dir, t)
}

func writeTemplate(dir string, tmpl interface{}) error {
	b, err := yaml.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("unable to marshal object: %v", err)
	}

	filename := filepath.Join(dir, "template.yaml")
	if err := os.WriteFile(filename, b, 0644); err != nil {
		return fmt.Errorf("unable to write template to file: %v", err)
	}

	return nil
}

// escapeName removes any characters from a name that are not valid in a filename
func escapeName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch r {
		// common set of reserved characters from
		// https://en.wikipedia.org/wiki/Filename#Comparison_of_filename_limitations
		// which will just be removed from the name
		case '|', '\\', '?', '*', '<', '>', '"', ':', '/':
			continue

		default:
			if !unicode.In(r, unicode.Cc) {
				b.WriteRune(r)
			}
		}
	}

	return b.String()
}
