// See LICENSE file for copyright and license details

package env

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"

	"djmo.ch/go-shrt/cmd/shrt/internal/base"
)

func TestRunEnv(t *testing.T) {
	var (
		w    = new(strings.Builder)
		args = make([]string, 0, 0)
		ctx  = context.Background()
	)
	ctx = context.WithValue(ctx, "w", w)
	ctx = context.WithValue(ctx, "args", args)
	clearEnv()
	runEnv(ctx)
	s := bufio.NewScanner(strings.NewReader(w.String()))
	for s.Scan() {
		kv := strings.SplitN(s.Text(), "=", 2)
		if len(kv) == 1 {
			t.Error("found key but no value")
		}
		if !strings.Contains(base.KnownEnv, kv[0]) {
			t.Error("found unexpected key ", kv[0])
		}
		if kv[1] != `""` {
			t.Error("unexpected value ", kv[1])
		}
	}
}

func TestRunEnvUW(t *testing.T) {
	envPath := "testdata/env"
	os.Setenv("SHRTENV", envPath)
	envMap := readEnvFile(envPath)
	if _, ok := envMap[base.SHRT_DBPATH]; !ok {
		t.Fatalf("%s key not found in %s", base.SHRT_DBPATH, envPath)
	}
	runEnvU([]string{base.SHRT_DBPATH})
	envMap = readEnvFile(envPath)
	if _, ok := envMap[base.SHRT_DBPATH]; ok {
		t.Fatalf("%s key found in %s", base.SHRT_DBPATH, envPath)
	}
	runEnvW([]string{base.SHRT_DBPATH + "=bar"})
	envMap = readEnvFile(envPath)
	if _, ok := envMap[base.SHRT_DBPATH]; !ok {
		t.Fatalf("%s key not found in %s", base.SHRT_DBPATH, envPath)
	}
}

func TestMergeEnv(t *testing.T) {
	clearEnv()
	envPath := "testdata/env"
	os.Setenv("DGITENV", envPath)
	MergeEnv()
	envMap := readEnvFile(envPath)
	if envMap[base.SHRT_SRVNAME] != "foo" {
		t.Errorf("unexpected value for %s: %s",
			base.SHRT_SRVNAME, envMap[base.SHRT_SRVNAME])
	}
}

func clearEnv() {
	for _, envVar := range strings.Fields(base.KnownEnv) {
		os.Unsetenv(envVar)
	}
}
