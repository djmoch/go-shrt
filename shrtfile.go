// See LICENSE file for copyright and license details

package shrt

import (
	"bufio"
	"fmt"
	"io/fs"
	"strings"
	"sync"
)

// ShrtType is the type of a ShrtFile entry. Their textual
// representations within the ShrtFile are listed in [NoneType].
type ShrtType int

const (
	NoneType  ShrtType = iota
	ShortLink          // shrtlnk
	GoGet              // goget
)

// ShrtEntry is a ShrtFile entry.
type ShrtEntry struct {
	URL  string
	Type ShrtType
}

// The ShrtFile struct contains the data read from a specially-formatted
// file. The syntax of the file is human readable. Each line
// represents a key-value pair. The key is everything to the left
// of the first equals sign, and the value is everything to the
// right. The value is then split around the first occurence of the
// colon character, with the left side representing the type, and the
// right side representing the URL. Whitespace is trimmed from the
// beginning and end of all fields.
type ShrtFile struct {
	m   map[string]ShrtEntry
	mux sync.RWMutex
}

// The NewShrtFile function returns a new ShrtFile.
func NewShrtFile() *ShrtFile {
	return &ShrtFile{m: make(map[string]ShrtEntry)}
}

// The ReadShrtFile function reads an existing ShrtFile from f and
// returns a pointer to a ShrtFile object. The provided file is closed
// before returning.
func (s *ShrtFile) ReadShrtFile(f fs.File) error {
	defer f.Close()
	s.mux.Lock()
	defer s.mux.Unlock()

	s.m = make(map[string]ShrtEntry)

	scnr := bufio.NewScanner(f)

	for scnr.Scan() {
		tok := strings.SplitN(scnr.Text(), "=", 2)
		if _, ok := s.m[strings.Trim(tok[0], " ")]; ok {
			return fmt.Errorf("repeat key: %s", tok[0])
		}
		if len(tok) != 2 {
			return fmt.Errorf("invalid syntax: %s", scnr.Text())
		}
		key := strings.TrimSpace(tok[0])
		tok = strings.SplitN(tok[1], ":", 2)
		if len(tok) != 2 {
			return fmt.Errorf("invalid syntax: %s", scnr.Text())
		}
		var typ ShrtType
		switch strings.TrimSpace(tok[0]) {
		case "shrtlnk":
			typ = ShortLink
		case "goget":
			typ = GoGet
		default:
			return fmt.Errorf("unrecognized type: %s", tok[0])
		}
		s.m[key] = ShrtEntry{
			Type: typ,
			URL:  strings.TrimSpace(tok[1]),
		}
	}
	return nil
}

// The Get method gets the value of the specified key. If the key
// does not exist, an error is returned.
func (s *ShrtFile) Get(key string) (ShrtEntry, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	entry, ok := s.m[key]
	if !ok {
		entry.Type = NoneType
		return entry, fmt.Errorf("key not found: %s", key)
	}
	return entry, nil
}
