package shrt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// The ShrtFile struct contains the data read from a specially-formatted
// file. The syntax of the file is human readable. Each line
// represents a key-value pair. The key is everything to the left
// of the first equals sign, and the value is everything to the
// right. Whitespace is trimmed from the beginning and end of both
// keys and values.
type ShrtFile struct {
	m    map[string]string
	path string
}

// The NewShrtFile function creates a new ShrtFile. The filesystem
// is checked for the existence of a node at the specified path,
// and the NewShrtFile returns an error if something is there. This
// constitutes a weak check, since a file could easily be created
// before Write() is called, in which case the existing file will
// be truncated.
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

// The ReadShrtFile function reads an existing ShrtFile from the
// filesystem and returns a pointer to a ShrtFile object. See the
// ShrtFile documentation for the expected file format.
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

// The Get method gets the value of the specified key. If the key
// does not exist, an empty string is returned.
func (s *ShrtFile) Get(key string) string {
	return s.m[key]
}

// The Put method adds a the specified value, associating it with
// the specified key. The value is overwritten if the key already
// exists.
func (s *ShrtFile) Put(key, value string) {
	s.m[key] = value
}

// The Write method serializes the contents of the ShrtFile to the
// file path specified when the ShrtFile was created. If a file
// already exists at the specified path, it is truncated and
// overwritten.
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
