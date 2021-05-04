package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func unite(args []string) {

}

func uniteTemplate(dir string, w io.Writer) error {
	kinds, err := listKindDirs(dir)
	if err != nil {
		return fmt.Errorf("unable to read dir %q: %w", dir, err)
	}

	fmt.Println("Kinds:", kinds)

	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()

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

			var obj interface{}
			var err error
			switch k {
			case kindDashboard:
				obj, err = readDashboard(dir)

			case kindTask:
				obj, err = readTask(dir)

			default:
				obj, err = readObject(dir)
			}
			if err != nil {
				return err
			}

			if err := enc.Encode(obj); err != nil {
				return fmt.Errorf("unable to encode object from %q: %v", dir, err)
			}
		}
	}

	return nil
}

func readDashboard(dir string) (*dashboard, error) {
	var db dashboard
	if err := readTemplate(dir, &db); err != nil {
		return nil, err
	}

	for i, c := range db.Spec.Charts {
		for j, q := range c.Queries {
			if !strings.HasPrefix(q.Query, "file://") {
				continue
			}

			filename := filepath.Join(dir, q.Query)
			b, err := ioutil.ReadFile(filepath.Join(dir, strings.TrimPrefix(q.Query, "file://")))
			if err != nil {
				return nil, fmt.Errorf("unable to read query from file %s: %w", filename, err)
			}

			db.Spec.Charts[i].Queries[j].Query = string(b)
		}
	}

	return &db, nil
}

func readTask(dir string) (*task, error) {
	var t task
	if err := readTemplate(dir, &t); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(t.Spec.Query, "file://") {
		return &t, nil
	}

	queryFile := filepath.Join(dir, strings.TrimPrefix(t.Spec.Query, "file://"))
	b, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read query from file %s: %w", queryFile, err)
	}

	t.Spec.Query = string(b)
	return &t, nil
}

func readObject(dir string) (*object, error) {
	var obj object
	if err := readTemplate(dir, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func readTemplate(dir string, v interface{}) error {
	fmt.Println("reading template", dir)

	templateFile := filepath.Join(dir, "template.yml")
	f, err := os.Open(templateFile)
	if err != nil {
		return fmt.Errorf("unable to open file %q: %v", templateFile, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(v); err != nil {
		return fmt.Errorf("unable to decode template %q: %v", templateFile, err)
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
		kindLabel:     3,
		kindTask:      2,
		kindDashboard: 1,
	}
	sort.Slice(kinds, func(i, j int) bool {
		return priority[kinds[i]] > priority[kinds[j]]
	})

	return kinds, nil
}
