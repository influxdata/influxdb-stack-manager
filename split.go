package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

const splitUsage = `
Split a template file into separate templates and flux queries.

Usage:
  influxdb-stack-manager split <src> <dest>

Where src is a yaml template file, and dest is a directory.
`

func split(args []string) {
	if len(args) != 2 {
		log.Fatal("Expected exactly two args")
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("couldn't open template file %q: %v", args[0], err)
	}
	defer f.Close()

	if err := splitTemplate(args[1], f); err != nil {
		log.Fatalf("couldn't split template: %v", err)
	}
}

// split the contents of the reader into seperate templates and
func splitTemplate(dir string, r io.Reader) error {
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

		dir := filepath.Join(dir, obj.Kind, escapeName(walkNode(&obj.Spec, "name").Value))
		// Check whether we have a name collision, just error if so.
		if _, ok := seenNames[dir]; ok {
			return fmt.Errorf("name collision detected: %q appears more than once", dir)
		}
		seenNames[dir] = struct{}{}

		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("unable to make directory %q: %v", dir, err)
		}

		var queryNodes []queryNode
		switch obj.Kind {
		case kindDashboard:
			queryNodes = walkDashboard(&obj.Spec)

		case kindTask:
			queryNodes = walkTask(&obj.Spec)
		}

		queryNames := map[string]int{}
		for _, qn := range queryNodes {
			var name string
			if n, ok := queryNames[name]; ok {
				name = fmt.Sprintf("%s_%d.flux", qn.Name, n)
			} else {
				name = fmt.Sprintf("%s.flux", qn.Name)
			}
			queryNames[name]++

			// Write out the query to file
			if err := os.WriteFile(filepath.Join(dir, name), []byte(qn.Node.Value), 0644); err != nil {
				return fmt.Errorf("unable to write query to file '%s/%s': %v", dir, name, err)
			}
			qn.Node.Value = fmt.Sprintf("file://%s", name)
			qn.Node.Style = yaml.FlowStyle
		}

		filename := filepath.Join(dir, "template.yml")
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("unable to create file %q: %v", filename, err)
		}
		defer f.Close()

		enc := yaml.NewEncoder(f)
		enc.SetIndent(2)
		if err := enc.Encode(obj); err != nil {
			return fmt.Errorf("unable to marshal object: %v", err)
		}
	}
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
