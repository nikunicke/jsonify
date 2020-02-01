package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/json"
	"unicode"
)

func before(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return value
	}
	return value[0:pos]
}

func after(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return value
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func nestify_values(a []map[string]interface{}, key string, delim string)  {
	for i := 0; i < len(a); i++ {
		var result []interface{}
		if _, ok := a[i][key]; !ok {
			continue
		}
		newValue := strings.Split(a[i][key].(string), delim)
		for _, item := range newValue {
			result = append(result, strings.Split(item, " | "))
		}
		a[i][key] = result
	}
}

func get_reverse_depends(a []map[string]interface{})	{
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
					if packages["Package"].(string) == before(depends_item, " (") {
						rev_depends = append(rev_depends, other_packages["Package"].(string))
					}
				}
			}
		}
		packages["Reverse_Depends"] = rev_depends
	}
}


type data_object map[string]interface{}

func (container data_object) NestifyKeyValues(key string, delim string) {
	if _, ok := container[key]; !ok {
		return
	}
	var result []interface{}
	newValue := strings.Split(container[key].(string), delim)
	for _, item := range newValue {
		result = append(result, strings.Split(item, " | "))
	}
	container[key] = result
}

func main()  {
	var value		string
	var key			string

	output, err := ioutil.ReadFile("/var/lib/dpkg/status")
	if err != nil {
		fmt.Println(err)
		return
	}
	packages := strings.Split(string(output), "\n\n")
	result := []map[string]interface{}{}
	for i := 0; i < len(packages) - 1; i++ {
		container := data_object{}
		for _, row := range (strings.Split(packages[i], "\n")) {
			if unicode.IsSpace(rune(row[0])) {
				value = value + row
				container[key] = value // after(value, ": ")
			} else {
				value = after(before(row, "\n"), ": ")
				key = before(after(row, "\n"), ":")
				container[key] = value
			}
		}
		container.NestifyKeyValues("Depends", ", ")
		result = append(result, container)
	}
	get_reverse_depends(result)
	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("/root/data/dpkg/dpkg.json", b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

