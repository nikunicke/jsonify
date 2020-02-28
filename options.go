package main

import (
	"flag"
	"os"
	"path"
	"log"
)

func parseArgs(target *args) {
	dir, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	target.path = (flag.String("path", "/var/lib/dpkg/status", "file to parse"))
	target.save = (flag.String("save", path.Dir(dir) + "/data/dpkg.json", "save path"))
	flag.Parse()
	target.tail = flag.Args()
}

func setArguments(optionalArgs []string) [][]string {
	if len(optionalArgs) == 0 {
		return nil
	}
	if len(optionalArgs) % 2 != 0 {
		log.Fatal("[INPUT ERROR] - Each argument should be followed by a delimiter")
	}
	setOptionalArgs := make([][]string, len(optionalArgs) / 2)
	j := 0
	for i := 0; i < len(optionalArgs); i += 2 {
		setOptionalArgs[j] = make([]string, 2)
		setOptionalArgs[j][0] = optionalArgs[i]
		setOptionalArgs[j][1] = optionalArgs[i + 1]
		j += 1
	}
	return (setOptionalArgs)
}
