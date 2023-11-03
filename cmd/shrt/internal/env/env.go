// See LICENSE file for copyright and license details

// Package env implements the "shrt env" command
package env

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"djmo.ch/go-shrt"
	"djmo.ch/go-shrt/cmd/shrt/internal/base"
)

const (
	srvNameDefault = "example.com"
	scmTypeDefault = "git"
	suffixDefault  = ".git"
	rdrNameDefault = "github.com/user"
	bareRdrDefault = "example.org"
)

var Cmd = &base.Command{
	Name:      "env",
	Usage:     "shrt env [-u] [-w] [var ...]",
	ShortHelp: "print Shrt environment information",
	LongHelp: `Env prints Shrt environment information.

By default env prints information as a shell script. If one or more
variable names is given as arguments, env prints the value of each
named variable on its own line.

The -u flag requires one or more arguments and unsets
the default setting for the named environment variables,
if one has been set with 'shrt env -w'.

The -w flag requires one or more arguments of the
form NAME=VALUE and changes the default settings
of the named environment variables to the given values. If the same
NAME is provided multiple times, the last one takes effect.

For more about environment variables, see 'shrt help environment'.
	`,
}

var (
	envU = Cmd.Flags.Bool("u", false, "")
	envW = Cmd.Flags.Bool("w", false, "")
)

func init() {
	// break init cycle
	Cmd.Run = runEnv
}

func runEnv(ctx context.Context) {
	var (
		w    = ctx.Value("w").(io.Writer)
		args = ctx.Value("args").([]string)
	)
	if *envU && *envW {
		log.Fatal("cannot use -w with -u")
	}

	if *envU {
		runEnvU(args)
		return
	}

	if *envW {
		runEnvW(args)
		return
	}

	// Environment is already merged
	if len(args) > 0 {
		for _, arg := range args {
			fmt.Fprintln(w, os.Getenv(arg))
		}
		return
	}
	for _, key := range strings.Fields(base.KnownEnv) {
		value := os.Getenv(key)
		fmt.Fprintf(w, "%s=\"%s\"\n", key, value)
	}
}

func runEnvU(args []string) {
	envPath := envOrDefault(base.SHRTENV, envDefault)
	curEnv := readEnvFile(envPath)

	for _, arg := range args {
		delete(curEnv, arg)
	}

	writeEnvFile(envPath, curEnv)
}

func runEnvW(args []string) {
	envToWrite := make(map[string]string)
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) == 1 {
			log.Fatal("malformed argument: ", arg)
		}
		if !strings.Contains(base.KnownEnv, kv[0]) {
			log.Fatal("unknown env variable: ", kv[0])
		}
		envToWrite[kv[0]] = kv[1]
	}

	envPath := envOrDefault(base.SHRTENV, envDefault)
	curEnv := readEnvFile(envPath)

	for k, v := range envToWrite {
		if k == base.SHRTENV {
			log.Println(base.SHRTENV, "can only be set using the OS environment")
			continue
		}
		if k == base.SHRT_DBPATH && !filepath.IsAbs(v) {
			log.Println(base.SHRT_DBPATH, "must be an absolute path ... ignoring")
			continue
		}
		curEnv[k] = v
	}

	writeEnvFile(envPath, curEnv)
}

func readEnvFile(path string) map[string]string {
	envMap := make(map[string]string)
	envFile, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("error reading %s: %s", path, err)
		}
		return envMap
	}

	s := bufio.NewScanner(bytes.NewReader(envFile))
	for s.Scan() {
		kv := strings.SplitN(s.Text(), "=", 2)
		if len(kv) == 1 {
			log.Fatalf("malformed line in %s: %s", path, s.Text())
		}

		if !strings.Contains(base.KnownEnv, kv[0]) {
			continue
		}
		envMap[kv[0]] = kv[1]
	}

	return envMap
}

func writeEnvFile(path string, envMap map[string]string) {
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		log.Fatalf("failed to create directory %s: %s", filepath.Dir(path), err)
	}

	envFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to open %s for writing: %s", path, err)
	}
	defer envFile.Close()

	for k, v := range envMap {
		fmt.Fprintf(envFile, "%s=%s\n", k, v)
	}
}

// ConfigFromEnv returns a Config object matching the current
// environment.
func ConfigFromEnv() shrt.Config {
	return shrt.Config{
		SrvName: envOrDefault(base.SHRT_SRVNAME, srvNameDefault),
		ScmType: envOrDefault(base.SHRT_SCMTYPE, scmTypeDefault),
		Suffix:  envOrDefault(base.SHRT_SUFFIX, suffixDefault),
		RdrName: envOrDefault(base.SHRT_RDRNAME, rdrNameDefault),
		BareRdr: envOrDefault(base.SHRT_BARERDR, bareRdrDefault),
		// Trim the leading / to satisfy fs.FS
		DbPath: strings.TrimPrefix(envOrDefault(base.SHRT_DBPATH, dbPathDefault), "/"),
	}
}

// MergeEnv merges the program's environment with that specified in
// SHRTENV. Values already specified in the environment take
// precedence.
func MergeEnv() {
	envPath := envOrDefault(base.SHRTENV, envDefault)
	envFile, err := os.ReadFile(envPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("error reading %s: %s", envPath, err)
		}
		envFile = []byte{}
	}

	// Read envfile into environment
	s := bufio.NewScanner(bytes.NewReader(envFile))
	for s.Scan() {
		kv := strings.SplitN(s.Text(), "=", 2)
		if len(kv) == 1 {
			log.Fatal("malformed line in SHRTENV: ", s.Text())
		}

		key := kv[0]
		if !strings.Contains(base.KnownEnv, key) {
			log.Fatal("unknown env var: ", key)
		}
		value := kv[1]

		if _, ok := os.LookupEnv(key); !ok {
			os.Setenv(key, value)
		}
	}

	defaults := map[string]string{
		base.SHRTENV:      envDefault,
		base.SHRT_SRVNAME: srvNameDefault,
		base.SHRT_SCMTYPE: scmTypeDefault,
		base.SHRT_SUFFIX:  suffixDefault,
		base.SHRT_RDRNAME: rdrNameDefault,
		base.SHRT_BARERDR: bareRdrDefault,
		base.SHRT_DBPATH:  dbPathDefault,
	}

	// Populate missing environment variables with defaults
	for _, key := range strings.Fields(base.KnownEnv) {
		if _, ok := os.LookupEnv(key); !ok {
			os.Setenv(key, envOrDefault(key, defaults[key]))
		}
	}
}

func envOrDefault(key, d string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		return d
	}
	return env
}
