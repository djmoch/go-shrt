package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	errNone = iota
	errArgNum
	errDatabase
	errNoFile
	errRepeatToken
)

var (
	arg0 = path.Base(os.Args[0])

	shrt, cfg *shrtFile
)

func usage(r int) {
	fmt.Printf("usage: %s [-d dbpath] [-c cfgpath] [init]\n", arg0)
	os.Exit(r)
}

func handl(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}

	key := req.URL.Path
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}

	if req.URL.Query().Get("go-get") == "1" {
		repo := key
		resp := "<!DOCTYPE html>\n"
		resp += "<html>\n"
		resp += "<head>\n"
		resp += `<meta http-equiv="Content-Type" `
		resp += "content=\"text/html; charset=utf-8\"/>\n"
		resp += "<meta name=\"go-import\" "
		resp += fmt.Sprintf("content=\"%s/%s %s %s/%s%s\">\n",
			(*cfg)["srvname"], repo, (*cfg)["scmtype"],
			(*cfg)["rdrname"], repo, (*cfg)["suffix"])
		resp += `<meta http-equiv="refresh" content="0; `
		resp += fmt.Sprintf("url=https://godoc.org/%s/%s\">\n",
			(*cfg)["srvname"], repo)
		resp += "</head>\n"
		resp += "<body>\n"
		resp += `Redirecting to docs at <a href="https://godoc.org/`
		resp += fmt.Sprintf("%s/%s\">godoc.org/%s/%s</a>...\n",
			(*cfg)["srvname"], repo, (*cfg)["srvname"], repo)
		resp += "</body>\n"
		resp += "</html>\n"
		w.Write([]byte(resp))
		return
	}

	if val, ok := (*shrt)[key]; ok {
		w.Header().Add("Location", val)
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte("Redirecting\n"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Short link not found\n"))
	}

}

func main() {
	if len(os.Args) > 4 {
		usage(errArgNum)
	}

	dbpath := "shrt.db"
	cfgpath := "shrt.conf"
	//doinit := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h":
			usage(errNone)
		case "-d":
			i += 1
			dbpath = os.Args[i]
		case "-c":
			i += 1
			cfgpath = os.Args[i]
		case "init":
			//doinit = true
			break
		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option -- %s\n",
				arg0, os.Args[i])

		}
	}

	//	if doinit {
	//		doInit(dbpath)
	//		os.Exit(errNone)
	//	}

	cfg = readShrtFile(cfgpath)
	shrt = readShrtFile(dbpath)

	http.Handle("/", http.HandlerFunc(handl))
	http.ListenAndServe(":8080", nil)
}
