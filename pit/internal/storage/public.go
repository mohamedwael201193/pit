package storage

import (
	"fmt"
	"regexp"
	"strings"
)

// Public evidence takes a different path from private memory. Private memory is
// always encrypted with a workspace key. A public receipt is uploaded in the
// clear so that a third party holding only the storage root can download it and
// recompute the digest. The two paths are kept apart by explicit guards: a
// private upload must carry an encryption key, a public upload must not.

// PublicUploadArgs builds the official client argv for an unencrypted evidence
// upload. The payer key pays 0G gas for the log-entry transaction and is never
// the workspace memory key.
func PublicUploadArgs(j Job) ([]string, error) {
	if err := RejectUnofficialClient(j.CLI); err != nil {
		return nil, err
	}
	payer, err := NormalizePayerKey(j.PayerKey)
	if err != nil {
		return nil, err
	}
	if j.InputPath == "" || j.RPC == "" || j.Indexer == "" {
		return nil, fmt.Errorf("incomplete_upload")
	}
	if strings.TrimSpace(j.KeyHex) != "" {
		return nil, fmt.Errorf("public_upload_takes_no_encryption_key")
	}
	args := []string{
		"upload",
		"--url", j.RPC,
		"--indexer", j.Indexer,
		"--file", j.InputPath,
		"--key", payer,
		"--log-color-disabled",
	}
	if strings.TrimSpace(j.Flow) != "" {
		args = append(args, "--flow-address", j.Flow)
	}
	if err := PublicMustNotEncrypt(args); err != nil {
		return nil, err
	}
	return args, nil
}

// PublicDownloadArgs builds the official client argv for verifying a public
// root. The merkle proof flag is mandatory: without it the client would hand
// back bytes it never checked.
func PublicDownloadArgs(j Job) ([]string, error) {
	if err := RejectUnofficialClient(j.CLI); err != nil {
		return nil, err
	}
	if err := RejectBadRoot(j.Root); err != nil {
		return nil, err
	}
	if j.OutPath == "" || j.Indexer == "" {
		return nil, fmt.Errorf("incomplete_download")
	}
	if strings.TrimSpace(j.KeyHex) != "" {
		return nil, fmt.Errorf("public_download_takes_no_encryption_key")
	}
	args := []string{
		"download",
		"--indexer", j.Indexer,
		"--file", j.OutPath,
		"--root", j.Root,
		"--proof",
		"--log-color-disabled",
	}
	if err := DownloadMustProve(args); err != nil {
		return nil, err
	}
	return args, nil
}

// PublicMustNotEncrypt is the mirror of UploadMustEncrypt. It keeps a public
// evidence upload from silently acquiring a key and becoming unverifiable.
func PublicMustNotEncrypt(args []string) error {
	for _, a := range args {
		if a == "--encryption-key" || a == "--encrypt" {
			return fmt.Errorf("public_upload_takes_no_encryption_key")
		}
	}
	return nil
}

var hex32 = regexp.MustCompile(`0x[0-9a-fA-F]{64}`)

// ParseRoot pulls the merkle root the official client reports for the upload.
func ParseRoot(out []byte) string {
	return firstTagged(out, []string{"root = ", "root="}, "merkle root")
}

// ParseTx pulls the 0G chain transaction hash of the log-entry submission. This
// is the on-chain half of the evidence: the root is the content commitment, the
// transaction is when the chain recorded it.
func ParseTx(out []byte) string {
	return firstTagged(out, []string{"txhash=", "hash="}, "append log entry")
}

func firstTagged(out []byte, needles []string, prefer string) string {
	lines := strings.Split(string(out), "\n")
	if prefer != "" {
		for _, ln := range lines {
			if strings.Contains(strings.ToLower(ln), prefer) {
				for _, needle := range needles {
					if found := afterNeedle(ln, needle); found != "" {
						return found
					}
				}
			}
		}
	}
	for _, ln := range lines {
		for _, needle := range needles {
			if found := afterNeedle(ln, needle); found != "" {
				return found
			}
		}
	}
	return ""
}

func afterNeedle(line, needle string) string {
	low := strings.ToLower(line)
	idx := strings.Index(low, needle)
	if idx < 0 {
		return ""
	}
	return hex32.FindString(line[idx:])
}

// ProofValidated reports whether the official client said it checked the merkle
// proof for the bytes it wrote. Any other outcome is treated as a failure.
func ProofValidated(out []byte) bool {
	return strings.Contains(strings.ToLower(string(out)), "succeeded to validate the downloaded file")
}

// AlreadyStored reports whether the network already held these exact bytes. The
// root is still correct and still verifiable, but no new chain transaction was
// submitted, so the caller must not claim one.
func AlreadyStored(out []byte) bool {
	low := strings.ToLower(string(out))
	return strings.Contains(low, "already exists") || strings.Contains(low, "data already")
}
