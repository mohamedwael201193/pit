package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/evidence"
	"github.com/mohamedwael201193/pit/internal/proof"
	"github.com/mohamedwael201193/pit/internal/storage"
)

// cmdEvidence is the third-party verification surface. Someone who does not
// trust this desk can list what PIT published and re-check any root against the
// live 0G network: merkle proof from the official client, digest recomputed from
// the downloaded bytes, and a finalized copy confirmed by storage nodes the
// indexer names.
func cmdEvidence(args []string) {
	sub := ""
	if len(args) > 0 {
		sub = args[0]
	}
	switch sub {
	case "list", "":
		evidenceList()
	case "verify":
		evidenceVerify(args[1:])
	case "bind-payer":
		evidenceBindPayer(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "usage: pit evidence [list|verify --root 0x..|bind-payer]")
		os.Exit(2)
	}
}

func evidenceFiler() (proof.Filer, error) {
	st, err := cli.Load(stateDir())
	if err != nil {
		return proof.Filer{}, fmt.Errorf("evidence requires pit init first")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return proof.Filer{}, err
	}
	ch := config.For(net)
	key, err := evidence.PayerKey(stateDir(), st.Network, st.WorkspaceID)
	if err != nil {
		key = ""
	}
	return proof.Filer{
		CLI: storage.LookCLI(), RPC: ch.RPC, Indexer: ch.StorageIndexer, Flow: ch.StorageFlow,
		Explorer: ch.Explorer, Network: string(ch.Network), ChainID: ch.ChainID,
		PayerKey: key, Dir: filepath.Join(stateDir(), "proofs"),
	}, nil
}

func evidenceList() {
	f, err := evidenceFiler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	rows, err := proof.Index(f.Dir, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"receipts": rows, "count": len(rows), "network": f.Network,
			"chain_id": f.ChainID, "sign": false, "trade": false,
		})
		return
	}
	if len(rows) == 0 {
		fmt.Println("no evidence filed yet")
		fmt.Println("research or a posted order publishes a public receipt to 0G Storage")
		return
	}
	fmt.Println("filed evidence on", f.Network, "chain", f.ChainID)
	for _, row := range rows {
		fmt.Println()
		fmt.Println(row.FiledAt, row.Kind, row.Market, row.Verdict)
		fmt.Println("  root  ", row.Root)
		fmt.Println("  digest", row.Digest)
		if row.OID != "" {
			fmt.Println("  oid   ", row.OID)
		}
		if row.TxLink != "" {
			fmt.Println("  chain ", row.TxLink)
		}
	}
	fmt.Println()
	fmt.Println("re-check any row with: pit evidence verify --root <root>")
}

func evidenceVerify(args []string) {
	root := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--root" && i+1 < len(args) {
			root = strings.TrimSpace(args[i+1])
		}
	}
	if root == "" {
		fmt.Fprintln(os.Stderr, "usage: pit evidence verify --root 0x...")
		os.Exit(2)
	}
	f, err := evidenceFiler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	got, err := f.Verify(ctx, root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(got)
		if got.Failure != "" {
			os.Exit(1)
		}
		return
	}
	fmt.Println("root           ", got.Root)
	fmt.Println("kind           ", got.Kind)
	fmt.Println("merkle proof   ", boolWord(got.ProofValidated))
	fmt.Println("digest         ", got.Recomputed)
	fmt.Println("digest matches ", boolWord(got.DigestMatch))
	fmt.Println("public safe    ", boolWord(got.PublicSafe))
	fmt.Println("roles verified ", boolWord(got.RolesVerified))
	fmt.Printf("finalized on    %d of %d storage nodes\n", got.FinalizedNodes, len(got.Nodes))
	if got.TxLink != "" {
		fmt.Println("chain tx       ", got.TxLink)
		fmt.Println("root committed ", boolWord(got.AnchorBound))
		if got.Anchor != nil && got.Anchor.BlockNumber != "" {
			fmt.Println("anchor block   ", got.Anchor.BlockNumber, "flow", got.Anchor.Flow)
		}
	}
	if got.Failure != "" {
		fmt.Fprintln(os.Stderr, "FAILED:", got.Failure)
		os.Exit(1)
	}
	fmt.Println("VERIFIED against the live 0G network")
}

func boolWord(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func evidenceBindPayer(args []string) {
	raw := strings.TrimSpace(os.Getenv("PIT_OG_PAYER_KEY"))
	for i := 0; i < len(args); i++ {
		if args[i] == "--from-env" {
			continue
		}
	}
	if raw == "" {
		fmt.Fprintln(os.Stderr, "set PIT_OG_PAYER_KEY in the environment, then run pit evidence bind-payer")
		os.Exit(2)
	}
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "evidence requires pit init first")
		os.Exit(2)
	}
	if err := evidence.SavePayerKey(stateDir(), st.Network, st.WorkspaceID, raw); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("payer key stored in the OS keyring for", st.WorkspaceID)
	fmt.Println("it pays 0G gas for public receipts and can never sign a venue order")
}
