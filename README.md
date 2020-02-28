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
*	<b>As arguments are not supported yet, you need to remove or customize splitElements, splitValues and getReverseDepends to your needs when parsing other files than dpkg/status</b>.

## Flags
*	-h: Usage
*	-path: path to file which we wish to parse to JSON
*	-save: path of the file where we wish to save our new JSON file
```console
	./path_to_repo/jsonify -path=./var/lib/dpkg/status -save=./data/dpkg_status.json
```

## Arguments
Additional arguments are not supported right now.

Here you could for example name the keys which values should be nested objects or arrays. <b>Right now you have to customize splitValues or splitElements to do that kind of stuff</b>.

## Files
*	main.go			-	Designated start and end of the program
*	options.go		-	Parsing flags and additional arguments
*	parse.go		-	Parsing the file to JSON
*	jsonify			-	Executable
*	monitor_file.sh	-	For automated updates. Add to crontab if you wish.
```console
	(crontab -l 2>/dev/null; echo "* * * * * ~/scripts/monitor_file.sh /var/lib/dpkg/status") | crontab - 
```
