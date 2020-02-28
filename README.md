# Text to JSON

This program allows you to dynamically generate JSON arrays from structured text files.
By default, the program will read and parse /var/lib/dpkg/status to /path/of/executable/dpkg.json
To run the program with default settings on your system, simply execute the binary.
```console
	./path_to_repo/jsonify
```

If you wish to recompile, you need to have Golang installed
```console
	go build -o "binary_name"
```

The text files need to be of similar structure as the "test" file in this repo.
*	Each object should be separated by two newlines
*	A line should contain its respective key and value
	*	If a key cannot be found, the line will be appended to the value of the previous key. If there is no previous key, the line will be skipped.

## Flags
*	-h: Usage
*	-path: path to file which we wish to parse to JSON
*	-save: path of the file where we wish to save our new JSON file
```console
	./path_to_repo/jsonify -path=./var/lib/dpkg/status -save=./data/dpkg_status.json
```

## Arguments
Additional arguments are not supported right now.

Here you could for example name the keys which values should be nested objects or arrays.

## Files
*	main.go	-	Designated start and end of the program
*	options.go	-	Parsing flags and additional arguments
*	parse.go	-	Parsing the file to JSON