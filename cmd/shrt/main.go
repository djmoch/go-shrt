// See LICENSE file for copyright and license details

//go:generate go test -v -run=TestDocsUpToDate -fixdocs

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"djmo.ch/go-shrt/cmd/shrt/internal/base"
	"djmo.ch/go-shrt/cmd/shrt/internal/env"
	"djmo.ch/go-shrt/cmd/shrt/internal/help"
	"djmo.ch/go-shrt/cmd/shrt/internal/serve"
	"djmo.ch/go-shrt/cmd/shrt/internal/version"
)

func init() {
	base.Shrt.Subcommands = []*base.Command{
		serve.Cmd,
		env.Cmd,
		version.Cmd,

		help.EnvCmd,
	}
}

func usagefunc(r int) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "usage: %s [-hv] [-d dbpath] [-c cfgpath] [-l listenaddr] [init]\n",
			os.Args[0])
		os.Exit(r)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
	flag.Usage = usagefunc(0)
	flag.Parse()

	env.MergeEnv()
	cfg := env.ConfigFromEnv()

	args := flag.Args()
	if len(args) < 1 {
		usagefunc(1)()
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "args", args[1:])
	ctx = context.WithValue(ctx, "w", os.Stdout)
	ctx = context.WithValue(ctx, "cfg", cfg)

	if args[0] == "help" {
		help.Help(ctx)
		return
	}

	cmd := base.FindCommand(args[0])
	if cmd == nil {
		fmt.Fprintf(os.Stderr, "%s %s: unknown command\n", os.Args[0], args[0])
		fmt.Fprintf(os.Stderr, "Run '%s help' for usage\n", os.Args[0])
		os.Exit(1)
	}

	cmd.Flags.Parse(os.Args[2:])
	ctx = context.WithValue(ctx, "args", cmd.Flags.Args())

	cmd.Run(ctx)
}
