package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	goshrt "djmo.ch/go-shrt"
)

const (
	errNone = iota
	errArgNum
	errDatabase
	errShrtFile
	errRepeatToken
	errInit
)

var (
	arg0 = path.Base(os.Args[0])
	mux  = sync.RWMutex{}

	shrt, cfg *goshrt.ShrtFile
	version   string
)

func usage(r int) {
	fmt.Printf("usage: %s [-hv] [-d dbpath] [-c cfgpath] [-l listenaddr] [init]\n", arg0)
	os.Exit(r)
}

func print_version() {
	fmt.Printf("%s version %s\n", arg0, version)
	os.Exit(errNone)
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

	if strings.Contains(key, "/") {
		log.Println("bad request: " + key)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Request path not allowed\n"))
		return
	}

	if key == "robots.txt" {
		log.Println("incoming robot")
		resp := "# Welcome to Shrt\n"
		resp += "User-Agent: *\n"
		resp += "Disallow:\n"
		w.Write([]byte(resp))
		return
	}

	if key == "" && cfg.Get("barerdr") != "" {
		log.Println("shortlink request for /")
		w.Header().Add("Location", cfg.Get("barerdr"))
		w.WriteHeader(http.StatusFound)
		w.Write([]byte("Redirecting\n"))
		return
	}

	mux.RLock()
	defer mux.RUnlock()
	if val := shrt.Get(key); val != "" {
		log.Println("shortlink request for", key)
		w.Header().Add("Location", val)
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte("Redirecting\n"))
		return
	}

	repo := key
	log.Println("go-get request for", repo)
	resp := "<!DOCTYPE html>\n"
	resp += "<html>\n"
	resp += "<head>\n"
	resp += `<meta http-equiv="Content-Type" `
	resp += "content=\"text/html; charset=utf-8\"/>\n"
	resp += "<meta name=\"go-import\" "
	resp += fmt.Sprintf("content=\"%s/%s %s %s/%s%s\">\n",
		cfg.Get("srvname"), repo, cfg.Get("scmtype"),
		cfg.Get("rdrname"), repo, cfg.Get("suffix"))
	resp += `<meta http-equiv="refresh" content="0; `
	resp += fmt.Sprintf("url=https://godoc.org/%s/%s\">\n",
		cfg.Get("srvname"), repo)
	resp += "</head>\n"
	resp += "<body>\n"
	resp += `Redirecting to docs at <a href="https://godoc.org/`
	resp += fmt.Sprintf("%s/%s\">godoc.org/%s/%s</a>...\n",
		cfg.Get("srvname"), repo, cfg.Get("srvname"), repo)
	resp += "</body>\n"
	resp += "</html>\n"
	w.Write([]byte(resp))
	return
}

func doInit(path string) {
	r := bufio.NewReader(os.Stdin)
	m, err := goshrt.NewShrtFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error creating config -- %s\n", arg0, err.Error())
		os.Exit(errInit)
	}

	fmt.Printf("server name: ")
	srvname, err := r.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", arg0, err.Error())
		os.Exit(errInit)
	}
	m.Put("srvname", srvname)

	fmt.Printf("SCM type: ")
	scmtype, err := r.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", arg0, err.Error())
		os.Exit(errInit)
	}
	m.Put("scmtype", scmtype)

	fmt.Printf("repo suffix (blank for none): ")
	suffix, err := r.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", arg0, err.Error())
		os.Exit(errInit)
	}
	m.Put("suffix", suffix)

	fmt.Printf("redirect base url: ")
	rdrname, err := r.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", arg0, err.Error())
		os.Exit(errInit)
	}
	m.Put("rdrname", rdrname)

	fmt.Printf("bare redirect url: ")
	barerdr, err := r.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", arg0, err.Error())
		os.Exit(errInit)
	}
	m.Put("barerdr", barerdr)
	m.Write()
}

func hangup(dbpath string) {
	hup := make(chan os.Signal)
	signal.Notify(hup, syscall.SIGHUP)
	for {
		<-hup
		tmpShrt, err := goshrt.ReadShrtFile(dbpath)
		if err != nil {
			log.Printf("db error -- %s\n", err.Error())
		} else {
			mux.Lock()
			shrt = tmpShrt
			mux.Unlock()
		}
	}
}

func main() {
	var err error
	if len(os.Args) > 4 {
		usage(errArgNum)
	}

	dbpath := "shrt.db"
	cfgpath := "shrt.conf"
	listenaddr := ":8080"
	doinit := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h":
			usage(errNone)
		case "-v":
			print_version()
		case "-d":
			i += 1
			dbpath = os.Args[i]
		case "-c":
			i += 1
			cfgpath = os.Args[i]
		case "-l":
			i += 1
			listenaddr = os.Args[i]
		case "init":
			doinit = true
			break
		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option -- %s\n",
				arg0, os.Args[i])

		}
	}

	if doinit {
		doInit(cfgpath)
		os.Exit(errNone)
	}

	cfg, err = goshrt.ReadShrtFile(cfgpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: config error -- %s\n", arg0, err.Error())
		os.Exit(errShrtFile)
	}

	shrt, err = goshrt.ReadShrtFile(dbpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: db error -- %s\n", arg0, err.Error())
		os.Exit(errShrtFile)
	}

	go hangup(dbpath)

	http.Handle("/", http.HandlerFunc(handl))
	log.Println("listening on", listenaddr)
	log.Fatal(http.ListenAndServe(listenaddr, nil))
}
