package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const uniteUsage = `
Unite a set of templates/flux queries back into a single template.

Usage:
  influxdb-stack-manager unite <src> <dest>

Where src is a directory, and dest is a template file.
`

// unite separated template files and flux queries into a single template.
func unite(args []string) error {
	if len(args) != 2 {
		log.Println(uniteUsage)
		return errors.New("expected exactly two args")
	}

	f, err := os.Create(args[1])
	if err != nil {
		return fmt.Errorf("couldn't create template file %q: %v", args[1], err)
	}
	defer f.Close()

	if err := uniteTemplate(args[0], f); err != nil {
		return fmt.Errorf("couldn't split template: %v", err)
	}

	return nil
}

// uniteTemplate walks a directory, finding all templates, reintegrating any flux queries that have been
// separated into their own files, and then writing them back to the writer.
func uniteTemplate(dir string, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()

	kinds, err := listKindDirs(dir)
	if err != nil {
		return fmt.Errorf("unable to read dir %q: %w", dir, err)
	}

	for _, k := range kinds {
		dir := filepath.Join(dir, k)
		items, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("unable to read dir %q: %w", dir, err)
		}

		for _, item := range items {
			if !item.IsDir() {
				continue
			}
			dir := filepath.Join(dir, item.Name())

			tf := filepath.Join(dir, templateFile)
			b, err := os.ReadFile(tf)
			if err != nil {
				return fmt.Errorf("unable to read file %q: %v", tf, err)
			}

			var obj object
			if err := yaml.Unmarshal(b, &obj); err != nil {
				return fmt.Errorf("unable to decode template %q: %v", tf, err)
			}

			// Find all query strings that need to be reunited.
			var queryNodes []queryNode
			switch obj.Kind {
			case kindCheck:
				queryNodes = walkCheck(&obj.Spec)

			case kindDashboard:
				queryNodes = walkDashboard(&obj.Spec)

			case kindTask:
				queryNodes = walkTask(&obj.Spec)
			}

			for _, qn := range queryNodes {
				if !strings.HasPrefix(qn.Node.Value, queryPrefix) {
					continue
				}

				filename := filepath.Join(dir, strings.TrimPrefix(qn.Node.Value, queryPrefix))
				b, err := os.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("unable to read query file %q: %v", filename, err)
				}
				qn.Node.SetString(string(b))
			}

			if err := enc.Encode(obj); err != nil {
				return fmt.Errorf("unable to encode object: %v", err)
			}
		}
	}

	return nil
}

// listKindDirs returns a sorted list of directories
func listKindDirs(dir string) ([]string, error) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var kinds []string
	for _, d := range dirs {
		if d.IsDir() {
			kinds = append(kinds, d.Name())
		}
	}

	// order which kinds should be added to the combined template.
	// higher numbers are added first.
	priority := map[string]int{
		kindLabel:     4,
		kindCheck:     3,
		kindTask:      2,
		kindDashboard: 1,
	}
	sort.Slice(kinds, func(i, j int) bool {
		return priority[kinds[i]] > priority[kinds[j]]
	})

	return kinds, nil
}
