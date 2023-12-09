// See LICENSE file for copyright and license details

// Package base defines the foundational structures required to build
// out the Shrt command suite.
package base

import (
	"context"
	"flag"
)

// Environment variable keys
const (
	SHRTENV           = "SHRTENV"
	SHRT_SRVNAME      = "SHRT_SRVNAME"
	SHRT_SCMTYPE      = "SHRT_SCMTYPE"
	SHRT_SUFFIX       = "SHRT_SUFFIX"
	SHRT_RDRNAME      = "SHRT_RDRNAME"
	SHRT_BARERDR      = "SHRT_BARERDR"
	SHRT_DBPATH       = "SHRT_DBPATH"
	SHRT_GOSOURCEDIR  = "SHRT_GOSOURCEDIR"
	SHRT_GOSOURCEFILE = "SHRT_GOSOURCEFILE"
)

// KnownEnv is a list of environment variables that affect the
// operation of the shrt command
const KnownEnv = `
	SHRTENV
	SHRT_SRVNAME
	SHRT_SCMTYPE
	SHRT_SUFFIX
	SHRT_RDRNAME
	SHRT_BARERDR
	SHRT_DBPATH
	SHRT_GOSOURCEDIR
	SHRT_GOSOURCEFILE
	`

type Command struct {
	Run                              func(context.Context)
	Flags                            flag.FlagSet
	Name, ShortHelp, LongHelp, Usage string
	Subcommands                      []*Command
}

var Shrt = &Command{
	Name: "shrt",
	LongHelp: `Shrt is a URL shortener and go-get redirector.

Shrt is a URL shortener service (much like bit.ly without the
trackers) that also handles go-get requests. The latter are a
specific GET request query used by the Go programming language
toolchain to aid in the downloading of utilities and libraries
prior to build and installation.

Upon invocation, shrt does one of two things depending
on the presence or absence of the init argument. If the init
argument is present, a series of questions is asked, the responses
are recorded in a configuration file, and the program exits. If
the init argument is absent, shrt reads the configuration and
database files into memory, binds to the port specified by the
-l flag, and begins serving requests.

Shortlinks are recorded in the database, and any request path
not matching a shortlink is assumed to be a go-get request. This
is by design, but can result in specious redirects. Additionally,
subdirectory paths are not allowed.

Shortlinks generate an HTTP 301 response. Go-get requests generate
an HTTP 200 response. If configured, requests to the base path
(i.e., "/") generate an HTTP 302 response.

In order to add a new shortlink to the database, simply edit the
file. After saving, users on Unix systems may send SIGHUP to a
running server process to reload the file. Non-Unix users will need
to restart the server.
`,
	Usage: "shrt <command> [arguments]",
}

func FindCommand(cmd string) *Command {
	for _, sub := range Shrt.Subcommands {
		if sub.Name == cmd {
			return sub
		}
	}
	return nil
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}
