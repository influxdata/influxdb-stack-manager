package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const uniteUsage = `
Unite a set of templates/flux queries back into a single template.

Usage:
  influxdb-stack-manager unite <src> <dest> [flags]

Where src is a directory, and dest is a template file.

Flags:
`

// unite separated template files and flux queries into a single template.
func unite(args []string) error {
	var dataFile string
	var help bool
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	fs.StringVar(&dataFile, "data-file", "", "Data file to use for injected data in templates")
	fs.BoolVarP(&help, "help", "h", false, "Display help for this command.")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("Error: %v\nSee 'influxdb-stack-manager push -h' for help", err)
	}

	args = fs.Args()
	if help {
		log.Println(uniteUsage + fs.FlagUsages())
		return nil
	}
	if len(args) != 2 {
		return errors.New("Error: wrong number of args\nSee 'influxdb-stack-manager unite -h' for help")
	}

	f, err := os.Create(args[1])
	if err != nil {
		return fmt.Errorf("couldn't create template file %q: %v", args[1], err)
	}
	defer f.Close()

	if err := uniteTemplate(args[0], f, dataFile); err != nil {
		return fmt.Errorf("couldn't split template: %v", err)
	}

	return nil
}

// uniteTemplate walks a directory, finding all templates, reintegrating any flux queries that have been
// separated into their own files, and then writing them back to the writer.
func uniteTemplate(dir string, w io.Writer, dataFile string) error {
	data, err := loadDataFile(dataFile)
	if err != nil {
		return fmt.Errorf("unable to load data file: %v", err)
	}

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

			tmpl, err := template.ParseGlob(filepath.Join(dir, "*"))
			if err != nil {
				return fmt.Errorf("unable to parse files in %q: %v", dir, err)
			}
			tmpl = tmpl.Option("missingkey=error")

			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, templateFile, data); err != nil {
				return fmt.Errorf("unable to execute template file %q: %v", filepath.Join(dir, templateFile), err)
			}

			var obj object
			if err := yaml.Unmarshal(buf.Bytes(), &obj); err != nil {
				return fmt.Errorf("unable to decode template %q: %v", filepath.Join(dir, templateFile), err)
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

			if len(queryNodes) > 0 {
				tmpl, err := template.ParseGlob(filepath.Join(dir, "*.flux"))
				if err != nil {
					return fmt.Errorf("unable to parse query files in %q: %s", dir, err)
				}

				for _, qn := range queryNodes {
					if !strings.HasPrefix(qn.Node.Value, queryPrefix) {
						continue
					}

					filename := strings.TrimPrefix(qn.Node.Value, queryPrefix)
					var buf bytes.Buffer
					err := tmpl.ExecuteTemplate(&buf, filename, data)
					if err != nil {
						return fmt.Errorf("unable to execute query template %q: %v", filename, err)
					}

					qn.Node.SetString(buf.String())
				}
			}

			if err := enc.Encode(obj); err != nil {
				return fmt.Errorf("unable to encode object: %v", err)
			}
		}
	}

	return nil
}

func loadDataFile(filename string) (interface{}, error) {
	if filename == "" {
		return nil, nil
	}

	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading data file %q: %v", filename, err)
	}

	var data interface{}
	switch ext := filepath.Ext(filename); ext {
	case ".json":
		err = json.Unmarshal(f, &data)

	case ".yaml", ".yml":
		err = yaml.Unmarshal(f, &data)

	default:
		return nil, fmt.Errorf("unrecognised data format %q", ext)
	}
	return data, err
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
