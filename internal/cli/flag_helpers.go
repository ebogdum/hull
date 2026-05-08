package cli

import (
	"strings"

	hullerr "github.com/ebogdum/hull/internal/errors"
)

// parseLabelFlags parses repeated --labels k=v entries into a map.
// Returns nil for an empty input so an empty map is not stored on releases.
func parseLabelFlags(raw []string) (map[string]string, error) {
	if 0 == len(raw) {
		return nil, nil
	}
	out := make(map[string]string, len(raw))
	for _, entry := range raw {
		k, v, found := strings.Cut(entry, "=")
		if !found || "" == k {
			return nil, hullerr.NewErrorf(hullerr.ErrCLIFlag, "invalid --labels entry %q (expected key=value)", entry)
		}
		out[k] = v
	}
	return out, nil
}
