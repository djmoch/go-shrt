// See LICENSE file for copyright and license details

// Package shrt implements a simple (perhaps simplistic) URL
// shortener. It also handles go-get requests.
//
// Shortlinks are recorded in the database, and any request path not
// matching a shortlink is assumed to be a go-get request. This is by
// design, but can result in specious redirects. Additionally,
// subdirectory paths are not allowed.
//
// Shortlinks generate an HTTP 301 response. Go-get requests generate
// an HTTP 200 response. If configured, requests to the base path
// (i.e., "/") generate an HTTP 302 response.
//
// The database file is human-readable. See [Shrtfile] for the full
// specification.
package shrt

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

var robotstxt = `# Welcome to Shrt
User-Agent: *
Disallow:
`
var shrtrsp = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>"
<meta name="go-import" content="{{ .SrvName }}/{{ .Repo }} {{ .ScmType }} {{ .URL }}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/{{ .SrvName }}/{{ .DocPath }}">
</head>
<body>
Redirecting to docs at <a href="https://godoc.org/{{ .SrvName }}/{{ .DocPath }}">godoc.org/{{ .SrvName }}/{{ .DocPath }}</a>...
</body>
</html>
`

type shrtRequest struct {
	SrvName string
	Repo    string
	ScmType string
	URL     string
	DocPath string
}

// Config contains all of the global configuration for Shrt. All
// values except BareRdr and DbPath are used in the go-import meta tag
// values for go-get requests.
type Config struct {
	// Server name of the Shrt host
	SrvName string
	// SCM (or VCS) type
	ScmType string
	// SCM repository suffix, if required by repository host
	Suffix string
	// The server name of the repository host
	RdrName string
	// Where requests with an empty path should redirect
	BareRdr string
	// The path to the [ShrtFile]-formatted database file.
	DbPath string
}

// ShrtHandler is the core [http.Handler] for go-shrt.
type ShrtHandler struct {
	ShrtFile *ShrtFile
	Config   Config
	FS       fs.FS
}

// Handle implements the http.Handler interface.
func (s *ShrtHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	p := req.URL.Path
	p = strings.TrimPrefix(p, "/")

	if p == "robots.txt" {
		log.Println("incoming robot")
		fmt.Fprint(w, robotstxt)
		return
	}

	if p == "" && s.Config.BareRdr != "" {
		log.Println("shortlink request for /")
		w.Header().Add("Location", s.Config.BareRdr)
		w.WriteHeader(http.StatusFound)
		fmt.Fprintln(w, "Redirecting")
		return
	}

	key := strings.SplitN(p, "/", 2)[0]

	val, err := s.ShrtFile.Get(key)
	if err != nil {
		log.Println("not found:", key)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch val.Type {
	case ShortLink:
		if key != p {
			log.Println("path elements following shortlink:", p)
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		log.Println("shortlink request for", key)
		w.Header().Add("Location", val.URL)
		w.WriteHeader(http.StatusMovedPermanently)
		fmt.Fprintln(w, "Redirecting")
	case GoGet:
		log.Println("go-get request for", key)
		t := template.Must(template.New("shrt").Parse(shrtrsp))
		sReq := shrtRequest{
			SrvName: s.Config.SrvName,
			Repo:    key,
			ScmType: s.Config.ScmType,
			URL:     val.URL,
			DocPath: p,
		}
		if err := t.Execute(w, sReq); err != nil {
			log.Println("error executing template:", err)
		}
	}
}
