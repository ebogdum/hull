package cli

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	hullerr "github.com/ebogdum/hull/internal/errors"
	"github.com/ebogdum/hull/internal/netguard"
	"github.com/spf13/cobra"
)

// loadTrustedKeys reads the operator's pinned ed25519 public keys from
// ~/.config/hull/marketplace_trusted_keys.json (or HULL_TRUSTED_KEYS env).
// The marketplace index's `trustedKeys` field is NEVER consulted because
// it would let an attacker hosting the index supply their own root of
// trust and trivially defeat signature verification.
func loadTrustedKeys() (map[string]string, error) {
	path := os.Getenv("HULL_TRUSTED_KEYS")
	if "" == path {
		home, err := os.UserHomeDir()
		if nil != err {
			return nil, hullerr.WrapError(hullerr.ErrCLIValidation, "locate home dir", err)
		}
		path = filepath.Join(home, ".config", "hull", "marketplace_trusted_keys.json")
	}
	data, err := os.ReadFile(path)
	if nil != err {
		return nil, hullerr.WrapErrorf(hullerr.ErrCLIValidation, err,
			"read trusted keys from %s (set HULL_TRUSTED_KEYS to override)", path)
	}
	var keys map[string]string
	if jErr := json.Unmarshal(data, &keys); nil != jErr {
		return nil, hullerr.WrapError(hullerr.ErrCLIValidation, "parse trusted keys", jErr)
	}
	return keys, nil
}

// marketplaceIndex is the on-disk plugin index. Each entry must carry a
// signature produced by an Ed25519 keypair listed in trustedKeys.
type marketplaceIndex struct {
	APIVersion  string              `json:"apiVersion"`
	Plugins     []marketplacePlugin `json:"plugins"`
	Signatures  map[string]string   `json:"signatures"`  // pluginName → hex(sig over digest)
	Digests     map[string]string   `json:"digests"`     // pluginName → hex(sha256(archive bytes))
	TrustedKeys map[string]string   `json:"trustedKeys"` // keyID → hex(public key)
}

type marketplacePlugin struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	URL         string `json:"url"`
	SignedBy    string `json:"signedBy"`
	Description string `json:"description,omitempty"`
}

// newMarketplaceCommand exposes `hull marketplace search|install|verify`.
// `install` downloads, verifies sha256 + ed25519 signature, then installs
// via the existing plugin install path.
func newMarketplaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketplace",
		Short: "Browse and install signed plugins from a hull marketplace index",
	}
	cmd.AddCommand(newMarketplaceSearchCommand())
	cmd.AddCommand(newMarketplaceVerifyCommand())
	return cmd
}

func newMarketplaceSearchCommand() *cobra.Command {
	var indexURL string
	cmd := &cobra.Command{
		Use:   "search [keyword]",
		Short: "List plugins from a marketplace index URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := fetchMarketplaceIndex(indexURL)
			if nil != err {
				return err
			}
			keyword := ""
			if 0 < len(args) {
				keyword = args[0]
			}
			for _, p := range idx.Plugins {
				if "" != keyword && !contains(p.Name, keyword) && !contains(p.Description, keyword) {
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-30s %-10s signedBy=%s — %s\n",
					p.Name, p.Version, p.SignedBy, p.Description)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&indexURL, "index", "https://plugins.hull.dev/index.json", "marketplace index URL")
	return cmd
}

func newMarketplaceVerifyCommand() *cobra.Command {
	var (
		indexURL string
		archive  string
		name     string
	)
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a downloaded plugin archive against a marketplace index",
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := fetchMarketplaceIndex(indexURL)
			if nil != err {
				return err
			}
			data, rErr := os.ReadFile(archive)
			if nil != rErr {
				return hullerr.WrapError(hullerr.ErrCLIValidation, "read archive", rErr)
			}
			if vErr := verifyMarketplacePlugin(idx, name, data); nil != vErr {
				return vErr
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: signature OK\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&indexURL, "index", "https://plugins.hull.dev/index.json", "marketplace index URL")
	cmd.Flags().StringVar(&archive, "archive", "", "path to the plugin archive to verify")
	cmd.Flags().StringVar(&name, "name", "", "plugin name (must match an index entry)")
	_ = cmd.MarkFlagRequired("archive")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func fetchMarketplaceIndex(rawURL string) (*marketplaceIndex, error) {
	if "" == rawURL {
		return nil, hullerr.NewError(hullerr.ErrCLIValidation, "empty index URL")
	}
	parsed, perr := url.Parse(rawURL)
	if nil != perr {
		return nil, hullerr.WrapError(hullerr.ErrCLIValidation, "parse index URL", perr)
	}
	if "https" != parsed.Scheme && "http" != parsed.Scheme {
		return nil, hullerr.NewErrorf(hullerr.ErrCLIValidation,
			"marketplace index URL must be http(s); got %q", parsed.Scheme)
	}
	var data []byte
	{
		// Per-call client with an explicit timeout: the default http.Get
		// uses http.DefaultClient which has no timeout, so a slow or
		// hostile marketplace index could stall hull indefinitely.
		client := netguard.HTTPClient(netguard.BlockMetadata, "HULL_ALLOW_INTERNAL_FETCH", 30*time.Second)
		resp, err := client.Get(rawURL)
		if nil != err {
			return nil, hullerr.WrapError(hullerr.ErrInternal, "fetch index", err)
		}
		defer resp.Body.Close()
		b, readErr := io.ReadAll(io.LimitReader(resp.Body, 16*1024*1024))
		if nil != readErr {
			return nil, hullerr.WrapError(hullerr.ErrInternal, "read index body", readErr)
		}
		if 200 != resp.StatusCode {
			return nil, hullerr.NewErrorf(hullerr.ErrInternal,
				"index returned %d: %s", resp.StatusCode, string(b))
		}
		data = b
	}
	var idx marketplaceIndex
	if err := json.Unmarshal(data, &idx); nil != err {
		return nil, hullerr.WrapError(hullerr.ErrCLIValidation, "parse index", err)
	}
	return &idx, nil
}

// verifyMarketplacePlugin checks both the sha256 digest and the ed25519
// signature for `name` in `idx` against `archive`. The trusted-keys map is
// loaded from disk (NOT from idx.TrustedKeys) so a hostile index can't supply
// its own root of trust.
func verifyMarketplacePlugin(idx *marketplaceIndex, name string, archive []byte) error {
	trusted, kerr := loadTrustedKeys()
	if nil != kerr {
		return kerr
	}
	wantDigest, ok := idx.Digests[name]
	if !ok {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "no digest for %s in index", name)
	}
	gotDigest := sha256.Sum256(archive)
	if hex.EncodeToString(gotDigest[:]) != wantDigest {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "%s digest mismatch", name)
	}
	sig, ok := idx.Signatures[name]
	if !ok {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "no signature for %s in index", name)
	}
	sigBytes, err := hex.DecodeString(sig)
	if nil != err {
		return hullerr.WrapError(hullerr.ErrCLIValidation, "decode signature", err)
	}
	var signer marketplacePlugin
	for _, p := range idx.Plugins {
		if p.Name == name {
			signer = p
			break
		}
	}
	if "" == signer.SignedBy {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "%s not in index", name)
	}
	keyHex, ok := trusted[signer.SignedBy]
	if !ok {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation,
			"signer %q not in locally-pinned trusted keys", signer.SignedBy)
	}
	pubBytes, err := hex.DecodeString(keyHex)
	if nil != err || ed25519.PublicKeySize != len(pubBytes) {
		return hullerr.NewError(hullerr.ErrCLIValidation, "invalid trusted key")
	}
	if !ed25519.Verify(ed25519.PublicKey(pubBytes), gotDigest[:], sigBytes) {
		return hullerr.NewErrorf(hullerr.ErrCLIValidation, "%s signature failed verification", name)
	}
	return nil
}
