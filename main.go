package main

import (
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var getwd = os.Getwd

func findTfFiles(root string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".tf" {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func readTF(file string) ([]string, error) {
	// read TF file and return all vars
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`var.(\w+)`)
	allMatches := r.FindAllStringSubmatch(string(content), -1)

	matches := []string{}
	for _, m := range allMatches {
		matches = append(matches, m[1])
	}
	return matches, nil
}

func removeDups(vars []string) []string {
	set := map[string]struct{}{}
	new := []string{}
	for _, v := range vars {
		if _, ok := set[v]; !ok {
			set[v] = struct{}{}
			new = append(new, v)
		}
	}
	return new
}

func makeVarsFile(vars []string, f io.Writer) error {
	// template vars into a vars tf file
	t := template.Must(template.New("vars").Parse(`{{range .}}variable "{{.}}" {}
{{end}}`))
	if err := t.Execute(f, vars); err != nil {
		return err
	}
	return nil
}

func main() {
	vars := []string{}
	path, err := getwd()
	if err != nil {
		log.Fatal(err)
	}
	files, err := findTfFiles(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		s, _ := readTF(file)
		vars = append(vars, s...)
	}
	vars = removeDups(vars)
	f, err := os.Create(filepath.Join(path, "vars.tf"))
    if err != nil {
        log.Fatal(err)
	}
    defer f.Close()

	if err = makeVarsFile(vars, f); err != nil {
		log.Fatal(err)
	}
}
