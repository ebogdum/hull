package repo

import (
	"crypto/subtle"
	"os"

	hullerr "github.com/ebogdum/hull/internal/errors"
	"github.com/ebogdum/hull/internal/logger"
)

// VerifyDigest checks a file's SHA256 against an expected digest.
// On mismatch, the file is deleted and an ErrDigest error is returned.
func VerifyDigest(filePath, expectedDigest string) error {
	actual, err := fileDigest(filePath)
	if nil != err {
		return hullerr.WrapError(hullerr.ErrDigest, "failed to compute file digest", err)
	}

	if 1 == subtle.ConstantTimeCompare([]byte(actual), []byte(expectedDigest)) {
		logger.Debug("digest verified for %s", filePath)
		return nil
	}

	logger.Warn("digest mismatch for %s: expected %s, got %s", filePath, expectedDigest, actual)

	if removeErr := os.Remove(filePath); nil != removeErr {
		logger.Warn("failed to remove file with mismatched digest: %s", removeErr)
	}

	return hullerr.NewErrorf(hullerr.ErrDigest,
		"digest mismatch for %s: expected %s, got %s", filePath, expectedDigest, actual)
}
