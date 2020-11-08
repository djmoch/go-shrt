package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type shrtFile map[string]string

func readShrtFile(db string) *shrtFile {
	var newFile shrtFile
	newFile = make(map[string]string)
	f, err := os.Open(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error opening shrtfile -- %s\n",
			arg0, err.Error())
		os.Exit(errNoFile)
	}

	defer f.Close()

	scnr := bufio.NewScanner(f)

	for scnr.Scan() {
		tok := strings.SplitN(scnr.Text(), "=", 2)
		if _, ok := newFile[strings.Trim(tok[0], " ")]; ok {
			fmt.Fprintf(os.Stderr, "%s: repeat token -- %s\n", arg0, tok[0])
			os.Exit(errRepeatToken)
		}
		newFile[strings.Trim(tok[0], " ")] = strings.Trim(tok[1], " ")
	}
	return &newFile
}
