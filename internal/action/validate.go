package action

import (
	"regexp"

	hullerr "github.com/ebogdum/hull/internal/errors"
)

var releaseNameRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

const maxReleaseNameLen = 53

// ValidateReleaseName checks that a release name is a valid DNS-1123 label subset.
// It returns an error if the name is empty, too long, or contains invalid characters.
func ValidateReleaseName(name string) error {
	if "" == name {
		return hullerr.NewError(hullerr.ErrCLIValidation, "release name is required")
	}

	if len(name) > maxReleaseNameLen {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "release name %q exceeds maximum length of %d characters", name, maxReleaseNameLen)
	}

	if !releaseNameRegex.MatchString(name) {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "release name %q is invalid: must match [a-z0-9]([a-z0-9-]*[a-z0-9])?", name)
	}

	return nil
}
