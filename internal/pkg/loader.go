package pkg

import (
	"os"
	"path/filepath"
	"regexp"

	hullerr "github.com/ebogdum/hull/internal/errors"
	"gopkg.in/yaml.v3"
)

var scopedNameRegex = regexp.MustCompile(`^(@[a-z0-9][a-z0-9-]*/)?[a-z0-9][a-z0-9.-]*$`)

const (
	packageFileName = "hull.yaml"
	valuesFileName  = "values.yaml"
)

// LoadMetadata is an alias for LoadPackageMetadata for concise usage.
func LoadMetadata(dirPath string) (PackageMetadata, error) {
	return LoadPackageMetadata(dirPath)
}

// LoadPackageMetadata reads and parses hull.yaml from the given directory path.
func LoadPackageMetadata(dirPath string) (PackageMetadata, error) {
	fullPath := filepath.Join(dirPath, packageFileName)
	data, err := os.ReadFile(fullPath)
	if nil != err {
		return PackageMetadata{}, hullerr.PackageError("failed to read hull.yaml", fullPath, err)
	}

	var meta PackageMetadata
	if err := yaml.Unmarshal(data, &meta); nil != err {
		return PackageMetadata{}, hullerr.PackageError("failed to parse hull.yaml", fullPath, err)
	}

	if validationErr := validateMetadata(&meta, fullPath); nil != validationErr {
		return PackageMetadata{}, validationErr
	}

	return meta, nil
}

// LoadValues reads and parses values.yaml from the given directory path.
func LoadValues(dirPath string) (Values, error) {
	fullPath := filepath.Join(dirPath, valuesFileName)
	data, err := os.ReadFile(fullPath)
	if nil != err {
		return nil, hullerr.PackageError("failed to read values.yaml", fullPath, err)
	}

	var vals Values
	if err := yaml.Unmarshal(data, &vals); nil != err {
		return nil, hullerr.PackageError("failed to parse values.yaml", fullPath, err)
	}

	return vals, nil
}

func validateMetadata(meta *PackageMetadata, filePath string) *hullerr.HullError {
	if "" == meta.APIVersion {
		return hullerr.NewError(hullerr.ErrPackageInvalid, "apiVersion is required").
			WithFile(filePath, 0, 0)
	}
	if "" == meta.Name {
		return hullerr.NewError(hullerr.ErrPackageInvalid, "name is required").
			WithFile(filePath, 0, 0)
	}
	if !scopedNameRegex.MatchString(meta.Name) {
		return hullerr.NewErrorf(hullerr.ErrPackageInvalid, "invalid package name %q: must match %s", meta.Name, scopedNameRegex.String()).
			WithFile(filePath, 0, 0)
	}
	if "" == meta.Version {
		return hullerr.NewError(hullerr.ErrPackageInvalid, "version is required").
			WithFile(filePath, 0, 0)
	}
	return nil
}
