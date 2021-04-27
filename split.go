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
Warning: This is a destructive operation, the destination directory will be
cleared if it already exists.
`

// split a file into separate templates and extract any flux code into its own file.
func split(args []string) error {
	if len(args) != 2 {
		log.Println(splitUsage)
		return errors.New("expected exactly two args")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("couldn't open template file %q: %v", args[0], err)
	}
	defer f.Close()

	if err := splitTemplate(args[1], f); err != nil {
		return fmt.Errorf("couldn't split template: %v", err)
	}

	return nil
}

// split the contents of the reader into separate templates and extract any flux code
// into their own files, organised under the supplied directory.
func splitTemplate(dir string, r io.Reader) error {
	// Clear the template directory, so we aren't left with any orphans.
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("unable to clear template directory: %w", err)
	}

	// Set of names that we've already seen, used to check for name collisions.
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
		if _, ok := seenNames[dir]; ok {
			// If we have a name collision, just error, we don't know how to organise two
			// objects with the same type and name.
			return fmt.Errorf("name collision detected: %q appears more than once", dir)
		}
		seenNames[dir] = struct{}{}

		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("unable to make directory %q: %v", dir, err)
		}

		// Find all of the query nodes present in the template.
		var queryNodes []queryNode
		switch obj.Kind {
		case kindDashboard:
			queryNodes = walkDashboard(&obj.Spec)

		case kindTask:
			queryNodes = walkTask(&obj.Spec)

		case kindCheck:
			queryNodes = walkCheck(&obj.Spec)
		}

		queryNames := map[string]int{}
		for _, qn := range queryNodes {
			// Keep track of used query names and, for any duplicates,
			// add a numerical suffix to distinguish them.
			var name string
			if n, ok := queryNames[name]; ok {
				name = fmt.Sprintf("%s_%d.flux", qn.Name, n)
			} else {
				name = fmt.Sprintf("%s.flux", qn.Name)
			}
			name = escapeName(name)
			queryNames[name]++

			// Write out the query to file
			filename := filepath.Join(dir, name)
			if err := os.WriteFile(filename, []byte(qn.Node.Value), 0644); err != nil {
				return fmt.Errorf("unable to write query to file %q: %v", filename, err)
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
	var previous rune
	for _, r := range name {
		switch r {
		// common set of reserved characters from
		// https://en.wikipedia.org/wiki/Filename#Comparison_of_filename_limitations
		// which will just be removed from the name
		case '|', '\\', '?', '*', '<', '>', '"', ':', '/':
			continue

		// Deduplicate multiple spaces caused by removal of characters
		case ' ':
			if previous == ' ' {
				continue
			}
			previous = r
			b.WriteRune(r)

		default:
			if !unicode.In(r, unicode.Cc) {
				previous = r
				b.WriteRune(r)
			}
		}
	}

	return b.String()
}
