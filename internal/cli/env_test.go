package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestEnvCommand_PrintsExpectedKeys(t *testing.T) {
	t.Setenv("HULL_NAMESPACE", "demo")
	cmd := newEnvCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.RunE(cmd, nil); nil != err {
		t.Fatalf("env: %v", err)
	}
	out := buf.String()
	for _, key := range []string{
		"HULL_BIN=", "HULL_CACHE_HOME=", "HULL_NAMESPACE=", "HULL_PLUGINS=",
	} {
		if !strings.Contains(out, key) {
			t.Errorf("env output missing %q\nfull:\n%s", key, out)
		}
	}
}

func TestEnvCommand_NamespaceFromLegacyEnvFallback(t *testing.T) {
	t.Setenv("HULL_NAMESPACE", "")
	t.Setenv("HELM_NAMESPACE", "fallback")
	env := collectHullEnv()
	if "fallback" != env["HULL_NAMESPACE"] {
		t.Errorf("legacy env-var fallback not honoured: %q", env["HULL_NAMESPACE"])
	}
}
