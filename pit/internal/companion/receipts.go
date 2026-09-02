package companion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/evidence"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/proof"
	"github.com/mohamedwael201193/pit/internal/receipt"
	"github.com/mohamedwael201193/pit/internal/storage"
)

// Evidence filing runs beside the product loop, never inside it. A research
// verdict or a posted order is already true when it happens; publishing it to 0G
// Storage is a second step that can be slow, so it never blocks the reply the
// user is waiting for. If publishing fails, the activity record says so instead
// of pretending a root exists.

func (h *Hub) proofDir() string {
	return filepath.Join(h.Dir, "proofs")
}

func (h *Hub) filer() (proof.Filer, error) {
	st, err := cli.Load(h.Dir)
	if err != nil {
		return proof.Filer{}, fmt.Errorf("unbound")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return proof.Filer{}, err
	}
	ch := config.For(net)
	key, err := evidence.PayerKey(h.Dir, st.Network, st.WorkspaceID)
	if err != nil {
		return proof.Filer{}, err
	}
	f := proof.Filer{
		CLI: storage.LookCLI(), RPC: ch.RPC, Indexer: ch.StorageIndexer, Flow: ch.StorageFlow,
		Explorer: ch.Explorer, Network: string(ch.Network), ChainID: ch.ChainID,
		PayerKey: key, Dir: h.proofDir(),
	}
	return f, f.Ready()
}

func (h *Hub) policyHash() string {
	pol := cli.ActivePolicy(h.Dir)
	hash, err := pol.Hash()
	if err != nil {
		return ""
	}
	return hash
}

func rolesForReceipt(roles []map[string]any) []receipt.Role {
	out := make([]receipt.Role, 0, len(roles))
	for _, rm := range roles {
		role := strings.TrimSpace(fmtString(rm["role"]))
		if role == "" {
			continue
		}
		item := receipt.Role{
			Role:       role,
			VerifyE2EE: strings.TrimSpace(fmtString(rm["verify_e2ee"])),
			Signer:     strings.TrimSpace(fmtString(rm["pubkey_signer"])),
			TeeSigner:  strings.TrimSpace(fmtString(rm["teeSigner"])),
			Side:       strings.TrimSpace(fmtString(rm["proposed_side"])),
		}
		if v, ok := rm["survives"].(bool); ok {
			item.Survives = v
		}
		if v, ok := rm["kill"].(bool); ok {
			item.Kill = v
		}
		out = append(out, item)
	}
	return out
}

func (h *Hub) newReceipt(kind string) (receipt.Receipt, config.Chain, error) {
	st, err := cli.Load(h.Dir)
	if err != nil {
		return receipt.Receipt{}, config.Chain{}, fmt.Errorf("unbound")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return receipt.Receipt{}, config.Chain{}, err
	}
	ch := config.For(net)
	return receipt.New(kind, string(ch.Network), ch.ChainID, st.WorkspaceID), ch, nil
}

// fileAsync publishes one receipt and records the outcome as an activity event.
// The caller returns immediately; the upload and the chain submission happen on
// their own goroutine.
func (h *Hub) fileAsync(r receipt.Receipt, market, jobID, oid string) {
	go func() {
		f, err := h.filer()
		if err != nil {
			appendActivity(h.Dir, activityEvent{
				WorkspaceID: r.Workspace, Kind: "evidence.unavailable", Market: market,
				Action: "evidence", Status: "blocked", JobID: jobID, OID: oid, Reason: err.Error(),
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cancel()
		filed, err := f.File(ctx, r)
		if err != nil {
			appendActivity(h.Dir, activityEvent{
				WorkspaceID: r.Workspace, Kind: "evidence.failed", Market: market,
				Action: "evidence", Status: "failed", JobID: jobID, OID: oid, Reason: err.Error(),
			})
			return
		}
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: r.Workspace, Kind: "evidence.filed", Market: market,
			Action: "evidence", Status: r.Kind, JobID: jobID, OID: oid,
			PreviewHash: r.PreviewHash, Root: filed.Root, Tx: filed.Tx,
			TxLink: filed.TxLink, Digest: filed.Digest, Link: filed.TxLink,
		})
	}()
}

// FileResearch publishes the committee verdict for a finished research job.
func (h *Hub) fileResearch(jobID, market, verdict, deny, previewHash string, roles []map[string]any, model, provider string) {
	r, _, err := h.newReceipt(receipt.KindResearch)
	if err != nil {
		return
	}
	r.Venue = "hyperliquid"
	r.Market = market
	r.Verdict = verdict
	r.Deny = deny
	r.PreviewHash = previewHash
	r.PolicyHash = h.policyHash()
	r.JobID = jobID
	r.Model = model
	r.Provider = provider
	r.Roles = rolesForReceipt(roles)
	if len(r.Roles) > 0 {
		r.TeeSigner = r.Roles[0].TeeSigner
	}
	h.fileAsync(r, market, jobID, "")
}

// FileOrder publishes a real venue action after the exchange accepted it.
func (h *Hub) fileOrder(got cli.OrderResult, jobID string) {
	if !got.Posted || strings.TrimSpace(got.OID) == "" {
		return
	}
	r, _, err := h.newReceipt(receipt.KindOrder)
	if err != nil {
		return
	}
	r.Venue = "hyperliquid"
	r.Market = got.Market
	r.Side = got.Side
	r.Size = got.Sz
	r.PreviewHash = got.Hash
	r.PolicyHash = h.policyHash()
	r.OID = got.OID
	r.Cloid = got.Cloid
	r.OrderStatus = strings.TrimSpace(got.Status)
	if r.OrderStatus == "" {
		r.OrderStatus = "posted"
	}
	r.JobID = jobID
	h.fileAsync(r, got.Market, jobID, got.OID)
}

func (h *Hub) localProofs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	body := map[string]any{"sign": false, "trade": false, "ready": false}
	f, err := h.filer()
	if err != nil {
		body["blocked"] = err.Error()
	} else {
		body["ready"] = true
		body["indexer"] = f.Indexer
		body["explorer"] = f.Explorer
		body["network"] = f.Network
		body["chain_id"] = f.ChainID
	}
	rows, _ := proof.Index(h.proofDir(), 60)
	out := make([]any, 0, len(rows))
	for _, row := range rows {
		if row.TxLink == "" && f.Explorer != "" {
			row.TxLink = proof.TxLink(f.Explorer, row.Tx)
		}
		out = append(out, row)
	}
	body["receipts"] = out
	body["count"] = len(out)
	writeLocal(w, http.StatusOK, body)
}

func (h *Hub) localProofVerify(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Root string `json:"root"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	f, err := h.filer()
	if err != nil {
		writeLocal(w, http.StatusConflict, map[string]any{"ok": false, "error": err.Error(), "sign": false, "trade": false})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
	defer cancel()
	got, err := f.Verify(ctx, strings.TrimSpace(body.Root))
	if err != nil {
		writeLocal(w, http.StatusBadRequest, map[string]any{"ok": false, "error": err.Error(), "sign": false, "trade": false})
		return
	}
	raw, err := json.Marshal(got)
	if err != nil {
		writeLocal(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": "encode", "sign": false, "trade": false})
		return
	}
	var shaped map[string]any
	_ = json.Unmarshal(raw, &shaped)
	shaped["ok"] = got.Failure == ""
	shaped["sign"] = false
	shaped["trade"] = false
	if got.Failure == "" {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "evidence.verified",
			Action: "verify", Status: got.Kind, Root: got.Root, Tx: got.Tx,
			TxLink: got.TxLink, Digest: got.Recomputed, Link: got.TxLink,
		})
	}
	writeLocal(w, http.StatusOK, shaped)
}

// evidenceModels reports the sealed models used for research so the receipt can
// name them without the desktop guessing.
func researchModels() (string, string) {
	return compute.MainnetChat().Model, "0g-compute-direct"
}
