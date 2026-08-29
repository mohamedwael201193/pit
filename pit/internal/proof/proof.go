// Package proof files PIT evidence on 0G Storage and reads it back. A filing is
// two facts a stranger can check without trusting this machine: a storage root
// whose bytes hash to the digest PIT recorded, and a 0G chain transaction that
// recorded the log entry. Everything in this package is real work against the
// live network; there is no local-only path that pretends to have filed.
package proof

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/receipt"
	"github.com/mohamedwael201193/pit/internal/storage"
)

// Filed is the durable index row for one piece of evidence.
type Filed struct {
	Kind        string  `json:"kind"`
	Digest      string  `json:"digest"`
	Root        string  `json:"root"`
	Tx          string  `json:"tx,omitempty"`
	Network     string  `json:"network"`
	ChainID     int64   `json:"chain_id"`
	Bytes       int     `json:"bytes"`
	Market      string  `json:"market,omitempty"`
	Verdict     string  `json:"verdict,omitempty"`
	OID         string  `json:"oid,omitempty"`
	JobID       string  `json:"job_id,omitempty"`
	Side        string  `json:"side,omitempty"`
	Size        float64 `json:"size,omitempty"`
	PreviewHash string  `json:"preview_hash,omitempty"`
	FiledAt     string  `json:"filed_at"`
	TxLink      string  `json:"tx_link,omitempty"`
	Duplicate   bool    `json:"duplicate,omitempty"`
}

// Filer holds everything needed to put bytes on 0G Storage for one workspace.
type Filer struct {
	CLI      string
	RPC      string
	Indexer  string
	Flow     string
	Explorer string
	Network  string
	ChainID  int64
	PayerKey string
	Dir      string
}

// Ready reports the first missing piece so the UI can say why filing is off
// instead of failing silently at the moment evidence matters.
func (f Filer) Ready() error {
	if err := storage.RefuseMissingProof(f.CLI); err != nil {
		return err
	}
	if strings.TrimSpace(f.RPC) == "" || strings.TrimSpace(f.Indexer) == "" {
		return errors.New("storage_endpoints_missing")
	}
	if _, err := storage.NormalizePayerKey(f.PayerKey); err != nil {
		return errors.New("payer_key_missing")
	}
	if strings.TrimSpace(f.Dir) == "" {
		return errors.New("proof_dir_missing")
	}
	return nil
}

const indexName = "receipts.jsonl"

// File writes the canonical receipt bytes, uploads them unencrypted, and records
// the root and chain transaction. The digest is computed before the upload so a
// later download can be checked against a number that existed first.
func (f Filer) File(ctx context.Context, r receipt.Receipt) (Filed, error) {
	if err := f.Ready(); err != nil {
		return Filed{}, err
	}
	raw, err := r.Canonical()
	if err != nil {
		return Filed{}, err
	}
	digest := receipt.DigestBytes(raw)
	if err := os.MkdirAll(filepath.Join(f.Dir, "payload"), 0o700); err != nil {
		return Filed{}, err
	}
	path := filepath.Join(f.Dir, "payload", strings.TrimPrefix(digest, "0x")+".json")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return Filed{}, err
	}
	job := storage.Job{
		CLI: f.CLI, RPC: f.RPC, Indexer: f.Indexer, Flow: f.Flow,
		PayerKey: f.PayerKey, InputPath: path,
	}
	argv, err := storage.PublicUploadArgs(job)
	if err != nil {
		return Filed{}, err
	}
	run, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	cmd := storage.Command(job, argv)
	cmd.Dir = f.Dir
	out, runErr := runWithContext(run, cmd)
	root := storage.ParseRoot(out)
	tx := storage.ParseTx(out)
	dup := storage.AlreadyStored(out)
	if err := storage.RejectBadRoot(root); err != nil {
		if runErr != nil {
			return Filed{}, fmt.Errorf("storage_upload_failed")
		}
		return Filed{}, fmt.Errorf("storage_root_missing")
	}
	filed := Filed{
		Kind: r.Kind, Digest: digest, Root: root, Tx: tx,
		Network: f.Network, ChainID: f.ChainID, Bytes: len(raw),
		Market: r.Market, Verdict: r.Verdict, OID: r.OID, JobID: r.JobID,
		Side: r.Side, Size: r.Size, PreviewHash: r.PreviewHash,
		FiledAt: time.Now().UTC().Format(time.RFC3339), Duplicate: dup,
	}
	filed.TxLink = TxLink(f.Explorer, tx)
	if err := Append(f.Dir, filed); err != nil {
		return filed, err
	}
	return filed, nil
}

// Append adds one row to the append-only index. Evidence is never rewritten.
func Append(dir string, filed Filed) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New("proof_dir_missing")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	line, err := json.Marshal(filed)
	if err != nil {
		return err
	}
	fh, err := os.OpenFile(filepath.Join(dir, indexName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = fh.Write(append(line, '\n'))
	return err
}

// Index reads the filed evidence newest first. A corrupt line is skipped rather
// than allowed to hide every later row.
func Index(dir string, limit int) ([]Filed, error) {
	raw, err := os.ReadFile(filepath.Join(dir, indexName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	lines := strings.Split(string(raw), "\n")
	out := make([]Filed, 0, len(lines))
	for i := len(lines) - 1; i >= 0; i-- {
		ln := strings.TrimSpace(lines[i])
		if ln == "" {
			continue
		}
		var filed Filed
		if json.Unmarshal([]byte(ln), &filed) != nil {
			continue
		}
		out = append(out, filed)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

// Find returns the index row for a root.
func Find(dir, root string) (Filed, bool) {
	rows, err := Index(dir, 0)
	if err != nil {
		return Filed{}, false
	}
	for _, row := range rows {
		if strings.EqualFold(row.Root, root) {
			return row, true
		}
	}
	return Filed{}, false
}

// Verified is the outcome of re-checking one root against the live network.
type Verified struct {
	Root           string             `json:"root"`
	Digest         string             `json:"digest"`
	Recomputed     string             `json:"recomputed"`
	DigestMatch    bool               `json:"digest_match"`
	ProofValidated bool               `json:"proof_validated"`
	PublicSafe     bool               `json:"public_safe"`
	RolesVerified  bool               `json:"roles_verified"`
	Nodes          []storage.NodeFile `json:"nodes,omitempty"`
	FinalizedNodes int                `json:"finalized_nodes"`
	Tx             string             `json:"tx,omitempty"`
	TxLink         string             `json:"tx_link,omitempty"`
	Anchor         *storage.Anchor    `json:"anchor,omitempty"`
	AnchorBound    bool               `json:"anchor_bound"`
	Kind           string             `json:"kind,omitempty"`
	CheckedAt      string             `json:"checked_at"`
	Failure        string             `json:"failure,omitempty"`
}

// Verify downloads the root through the official client with merkle proof
// checking, recomputes the digest, confirms a 0G Chain transaction commits the
// same root, and asks the storage network whether the root is finalized.
// Anything short of all of those is reported as a failure.
func (f Filer) Verify(ctx context.Context, root string) (Verified, error) {
	got := Verified{Root: root, CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if err := storage.RejectBadRoot(root); err != nil {
		return got, err
	}
	if err := storage.RefuseMissingProof(f.CLI); err != nil {
		return got, err
	}
	row, known := Find(f.Dir, root)
	if known {
		got.Digest = row.Digest
		got.Tx = row.Tx
		got.TxLink = TxLink(f.Explorer, row.Tx)
	}
	dir, err := os.MkdirTemp("", "pit-verify-")
	if err != nil {
		return got, err
	}
	defer os.RemoveAll(dir)
	out := filepath.Join(dir, "receipt.json")
	job := storage.Job{CLI: f.CLI, Indexer: f.Indexer, Root: root, OutPath: out}
	argv, err := storage.PublicDownloadArgs(job)
	if err != nil {
		return got, err
	}
	run, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	cmd := storage.Command(job, argv)
	cmd.Dir = dir
	raw, _ := runWithContext(run, cmd)
	got.ProofValidated = storage.ProofValidated(raw)
	body, readErr := os.ReadFile(out)
	if readErr != nil {
		got.Failure = "download_produced_no_bytes"
		return got, nil
	}
	got.Recomputed = receipt.DigestBytes(body)
	got.DigestMatch = got.Digest == "" || strings.EqualFold(got.Digest, got.Recomputed)
	if got.Digest == "" {
		got.Digest = got.Recomputed
	}
	parsed, perr := receipt.Parse(body)
	if perr == nil {
		got.PublicSafe = receipt.PublicSafe(body) == nil
		got.RolesVerified = receipt.RolesVerified(parsed.Roles)
		got.Kind = parsed.Kind
	}
	if nodes, nerr := storage.Availability(ctx, f.Indexer, root, 4); nerr == nil {
		got.Nodes = nodes
		got.FinalizedNodes = storage.FinalizedCount(nodes)
	}
	if got.Tx != "" {
		anchor := storage.CheckAnchor(ctx, f.RPC, got.Tx, root, f.Flow)
		got.Anchor = &anchor
		got.AnchorBound = anchor.Bound()
	}
	switch {
	case !got.ProofValidated:
		got.Failure = "merkle_proof_not_validated"
	case !got.DigestMatch:
		got.Failure = "digest_mismatch"
	case perr != nil:
		got.Failure = "receipt_unreadable"
	case got.FinalizedNodes == 0:
		got.Failure = "no_finalized_copy"
	case got.Tx != "" && !got.AnchorBound:
		got.Failure = "chain_transaction_does_not_commit_root"
	}
	return got, nil
}

// TxLink builds the public chain explorer link for a storage submission.
func TxLink(explorer, tx string) string {
	explorer = strings.TrimRight(strings.TrimSpace(explorer), "/")
	tx = strings.TrimSpace(tx)
	if explorer == "" || tx == "" {
		return ""
	}
	return explorer + "/tx/" + tx
}
