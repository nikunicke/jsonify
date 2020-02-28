package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

type args struct {
	path	*string
	save	*string
	tail	[]string
}

type packageObject map[string]interface{}
type packageArray []map[string]interface{}

func	main() {
	var target args
	var optionalArgs [][]string
	var file []byte
	var jsonResult packageArray
	
	parseArgs(&target)
	optionalArgs = setArguments(target.tail)
	file = validateFile(*target.path)
	jsonResult = parseFile(string(file), optionalArgs)
	getReverseDepends(jsonResult);
	jsonBytes, err := json.MarshalIndent(jsonResult, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*target.save, jsonBytes, 0644)
}
