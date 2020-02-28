package main

import (
	// "fmt"
	"flag"
	"os"
	"path"
	"log"
	"io/ioutil"
	"strings"
	"unicode"
	"encoding/json"
)

type cmd_args struct {
	path	*string
	save	*string
	tail	[]string
}

type data_object map[string]interface{}
type data_array []map[string]interface{}

func parse_cmd_flags(target *cmd_args) {
	dir, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	target.path = (flag.String("path", "/var/lib/dpkg/status", "file to parse"))
	target.save = (flag.String("save", path.Dir(dir) + "/data/dpkg.json", "save path"))
	flag.Parse()
	target.tail = flag.Args()
}

func validate_file(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return (file)
}

func set_arguments(args []string) [][]string {
	if len(args) == 0 {
		return nil
	}
	if len(args) % 2 != 0 {
		log.Fatal("[INPUT ERROR] - Each argument should be followed by a delimiter")
	}
	set_args := make([][]string, len(args) / 2)
	j := 0
	for i := 0; i < len(args); i += 2 {
		set_args[j] = make([]string, 2)
		set_args[j][0] = args[i]
		set_args[j][1] = args[i + 1]
		j += 1
	}
	return (set_args)
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

func get_reverse_depends(a data_array) {
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

func parse_file(file string, args [][]string) []map[string]interface{} {
	var value		string
	var key			string
	var json_array	data_array
	var row			int
	
	row = 0
	for _, pkg := range (strings.Split(file, "\n\n")) {
		if pkg == "" {
			continue
		}
		pkg_container := data_object{}
		for _, line := range (strings.Split(pkg, "\n")) {
			row++
			if !validate_line(line, row) {
				if key == "" { continue }
				value = value + "\n" + line
				pkg_container[key] = value
				} else {
					key = line[0:(strings.Index(line, ":"))]
					value = strings.TrimLeft(line[(strings.Index(line, ":") + 1):len(line)], "\t \n")
					pkg_container[key] = value
				}
			}
			pkg_container.SplitValues("Depends", ", ")
		pkg_container.SplitElements("Depends", " | ")
		json_array = append(json_array, pkg_container)
		row++
	}
	// fmt.Printf("%v\n\n", json_array)
	return (json_array)
}

func (container data_object) SplitValues(key string, delim string) {
	if _, ok := container[key]; !ok {
		return
	}
	array := strings.Split(container[key].(string), delim)
	container[key] = array;
}

func (container data_object) SplitElements(key string, delim string) {
	var result []interface{}
	
	if _, ok := container[key]; !ok {
		return
	}
	for _, element := range container[key].([]string) {
		result = append(result, strings.Split(element, delim))
	}
	container[key] = result
}

func	main() {
	var target cmd_args
	var additional_args [][]string
	var file []byte
	var json_result data_array
	
	parse_cmd_flags(&target)
	additional_args = set_arguments(target.tail)
	file = validate_file(*target.path)
	json_result = parse_file(string(file), additional_args)
	get_reverse_depends(json_result);
	json_bytes, err := json.MarshalIndent(json_result, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*target.save, json_bytes, 0644)
}
