// See LICENSE file for copyright and license details

// Package serve implements the "shrt serve" command
package serve

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"djmo.ch/go-shrt"
	"djmo.ch/go-shrt/cmd/shrt/internal/base"
)

var (
	hangup   func(*shrt.ShrtHandler)
	lockdown func(string)
)

var Cmd = &base.Command{
	Run:       runServe,
	Name:      "serve",
	Usage:     "shrt serve URL",
	ShortHelp: "serve requests",
	LongHelp: `Serve serves HTTP requests.

Shrt listens and serves shortlinks and go-get requests on the provided
URL. The only recognized scheme is http.
	`,
}

func runServe(ctx context.Context) {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("")
	var (
		args = ctx.Value("args").([]string)
		cfg  = ctx.Value("cfg").(shrt.Config)
	)
	if len(args) != 1 {
		log.Fatal("no URL provided")
	}
	u, err := url.Parse(args[0])
	if err != nil {
		log.Fatal("failed to parse URL: ", err)
	}

	shrtfile, err := shrt.ReadShrtFile(cfg.DbPath)
	if err != nil {
		log.Println("db error:", err)
		os.Exit(1)
	}
	h := &shrt.ShrtHandler{Config: cfg, ShrtFile: shrtfile}
	if hangup != nil {
		go hangup(h)
	}
	if lockdown != nil {
		lockdown(cfg.DbPath)
	}
	switch u.Scheme {
	case "http":
		listener, err := net.Listen("tcp", u.Host)
		if err != nil {
			log.Fatal("listen: ", err)
		}
		log.Fatal(http.Serve(listener, h))
	default:
		log.Fatal("unknown scheme:", u.Scheme)
	}
}
