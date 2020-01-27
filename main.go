package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/json"
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
	// fmt.Printf("%d\n",len(a))
	var result []interface{}
	for i := 0; i < len(a); i++ {
		if _, ok := a[i][key]; !ok {
			continue
		}
		newValue := strings.Split(a[i][key].(string), delim)
		for _, item := range newValue {
			result = append(result, item)
		}
		a[i][key] = result
	}
}

func main()  {
	var value	string
	var key		string

	output, err := ioutil.ReadFile("/var/lib/dpkg/status")
	if err != nil {
		fmt.Println(err)
		return
	}
	packages := strings.Split(string(output), "\n\n")
	result := []map[string]interface{}{}
	for i := 0; i < len(packages) - 1; i++ {
		container := map[string]interface{}{}
		for _, row := range (strings.Split(packages[i], "\n")) {
			if row[0] == ' ' {
				value = value + row
				container[key] = after(value, ": ")
			} else {
				value = after(before(row, "\n"), ": ")
				key = before(after(row, "\n"), ":")
				container[key] = value
			}
		}
		result = append(result, container)
	}
	nestify_values(result, "Depends", ", ")
	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("/home/npimenof/services/json_gen/result.json", b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

