package main

import (
	"log"
	"io/ioutil"
	"strings"
	"unicode"
)

func parseFile(file string, args [][]string) []map[string]interface{} {
	var value		string
	var key			string
	var pkgJson		packageArray
	var row			int
	
	row = 0
	for _, pkg := range (strings.Split(file, "\n\n")) {
		if pkg == "" {
			continue
		}
		pkgContainer := packageObject{}
		for _, line := range (strings.Split(pkg, "\n")) {
			row++
			if !validate_line(line, row) {
				if key == "" { continue }
				value = value + "\n" + line
				pkgContainer[key] = value
			} else {
					key = line[0:(strings.Index(line, ":"))]
					value = strings.TrimLeft(line[(strings.Index(line, ":") + 1):len(line)], "\t \n")
					pkgContainer[key] = value
			}
		}
		pkgContainer.SplitValues("Depends", ", ")
		pkgContainer.SplitElements("Depends", " | ")
		pkgJson = append(pkgJson, pkgContainer)
		row++
	}
	return (pkgJson)
}

func validateFile(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return (file)
}

func validate_line(line string, row int) bool {
	var index	int
	if index = strings.Index(line, ":"); index <= 0 {
		// fmt.Fprintln(os.Stderr, "[FORMAT WARNING]:", row,  "- Missing colon")
		return false
	}
	for _, char := range (line[0:index]) {
		if unicode.IsSpace(char) {
			// fmt.Fprintln(os.Stderr, "[FORMAT WARNING]:", row, "- Whitespace in object key")
			return false
		}
	}
	return true
}

func (container packageObject) SplitValues(key string, delim string) {
	if _, ok := container[key]; !ok {
		return
	}
	array := strings.Split(container[key].(string), delim)
	container[key] = array;
}

func (container packageObject) SplitElements(key string, delim string) {
	var result []interface{}
	
	if _, ok := container[key]; !ok {
		return
	}
	for _, element := range container[key].([]string) {
		result = append(result, strings.Split(element, delim))
	}
	container[key] = result
}

func getReverseDepends(a packageArray) {
	for _, packages := range a {
		var rev_depends []string
		for _, other_packages := range a {
			if packages["Package"].(string) == other_packages["Package"].(string) {
				continue
			}
			if _, ok := other_packages["Depends"]; !ok {
				continue
			}
			for _, depends := range other_packages["Depends"].([]interface{}) {
				for _, depends_item := range depends.([]string) {
					var length int;
					if length = strings.Index(depends_item, " ("); length == -1 {
						length = len(depends_item);
					}
					if packages["Package"].(string) == depends_item[0:length] {
						rev_depends = append(rev_depends, other_packages["Package"].(string))
					}
				}
			}
		}
		packages["Reverse_Depends"] = rev_depends
	}
}