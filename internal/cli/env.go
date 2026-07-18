package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

// newEnvCommand prints the resolved hull environment (paths, plugin dir, cache
// dir, namespace, kubeconfig).
func newEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Print hull's environment information",
		Long:  "Print resolved environment variables and paths used by hull.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := collectHullEnv()
			keys := make([]string, 0, len(env))
			for k := range env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%q\n", k, env[k])
			}
			return nil
		},
	}
	return cmd
}

func collectHullEnv() map[string]string {
	cacheRoot := os.Getenv("HULL_CACHE_HOME")
	if "" == cacheRoot {
		if home, err := os.UserCacheDir(); nil == err {
			cacheRoot = filepath.Join(home, "hull")
		}
	}
	configRoot := os.Getenv("HULL_CONFIG_HOME")
	if "" == configRoot {
		if home, err := os.UserConfigDir(); nil == err {
			configRoot = filepath.Join(home, "hull")
		}
	}
	dataRoot := os.Getenv("HULL_DATA_HOME")
	if "" == dataRoot {
		if home, err := os.UserHomeDir(); nil == err {
			dataRoot = filepath.Join(home, ".local", "share", "hull")
		}
	}
	pluginsRoot := os.Getenv("HULL_PLUGINS")
	if "" == pluginsRoot {
		pluginsRoot = filepath.Join(dataRoot, "plugins")
	}

	// Order: explicit --namespace flag > HULL_NAMESPACE > HELM_NAMESPACE > default.
	// HELM_NAMESPACE is honoured as a legacy fallback so existing shell
	// environments resolve to the same namespace every other hull command uses.
	ns := namespace
	if "" == ns {
		ns = os.Getenv("HULL_NAMESPACE")
	}
	if "" == ns {
		ns = os.Getenv("HELM_NAMESPACE")
	}
	if "" == ns {
		ns = "default"
	}

	bin, _ := os.Executable()
	if "" == bin {
		bin = os.Args[0]
	}

	return map[string]string{
		"HULL_BIN":               bin,
		"HULL_CACHE_HOME":        cacheRoot,
		"HULL_CONFIG_HOME":       configRoot,
		"HULL_DATA_HOME":         dataRoot,
		"HULL_PLUGINS":           pluginsRoot,
		"HULL_NAMESPACE":         ns,
		"HULL_KUBECONFIG":        os.Getenv("KUBECONFIG"),
		"HULL_KUBECONTEXT":       os.Getenv("HULL_KUBECONTEXT"),
		"HULL_REGISTRY_CONFIG":   filepath.Join(configRoot, "registry.json"),
		"HULL_REPOSITORY_CACHE":  filepath.Join(cacheRoot, "repository"),
		"HULL_REPOSITORY_CONFIG": filepath.Join(configRoot, "repositories.yaml"),
	}
}
