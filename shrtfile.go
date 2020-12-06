package shrt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ShrtFile struct {
	m    map[string]string
	path string
}

func NewShrtFile(path string) (*ShrtFile, error) {
	// we fail if the file already exists, so the logic is reversed
	// from the usual here
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		return nil,
			fmt.Errorf("File already exists. Please delete and try again.")
	}

	return &ShrtFile{m: make(map[string]string), path: path}, nil
}

func ReadShrtFile(db string) (*ShrtFile, error) {
	var newFile ShrtFile
	newFile.m = make(map[string]string)
	f, err := os.Open(db)
	if err != nil {
		return nil, fmt.Errorf("error opening shrtfile -- %s", err.Error())
	}

	defer f.Close()

	scnr := bufio.NewScanner(f)

	for scnr.Scan() {
		tok := strings.SplitN(scnr.Text(), "=", 2)
		if _, ok := newFile.m[strings.Trim(tok[0], " ")]; ok {
			return nil, fmt.Errorf("repeat token -- %s", tok[0])
		}
		newFile.m[strings.Trim(tok[0], " ")] = strings.Trim(tok[1], " ")
	}
	return &newFile, nil
}

func (s *ShrtFile) Get(key string) string {
	return s.m[key]
}

func (s *ShrtFile) Put(key, value string) {
	s.m[key] = value
}

func (s *ShrtFile) Write() error {
	f, err := os.Create(s.path)
	if err != nil {
		return fmt.Errorf("error opening ShrtFile: %s", err.Error())
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for k, v := range s.m {
		_, err = w.WriteString(fmt.Sprintf("%s = %s", k, v))
		if err != nil {
			return fmt.Errorf("error writing ShrtFile: %s", err.Error())
		}
	}
	w.Flush()
	return nil
}
