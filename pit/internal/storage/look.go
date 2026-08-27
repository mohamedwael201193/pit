package storage

import (
	"fmt"
	"os"
	"strings"
)

func LookCLI() string {
	return strings.TrimSpace(os.Getenv("PIT_STORAGE_CLI"))
}

func ProofJob(cliPath, rpc, indexer, keyHex, root, out string) (Job, error) {
	j := Job{CLI: cliPath, RPC: rpc, Indexer: indexer, KeyHex: keyHex, Root: root, OutPath: out}
	if _, err := DownloadArgs(j); err != nil {
		return Job{}, err
	}
	return j, nil
}

func RedactArgs(args []string) []string {
	out := make([]string, len(args))
	copy(out, args)
	for i := 0; i < len(out)-1; i++ {
		switch out[i] {
		case "--encryption-key", "--key", "--private-key":
			out[i+1] = "[redacted]"
		}
	}
	return out
}

func RefuseMissingProof(cliPath string) error {
	if strings.TrimSpace(cliPath) == "" {
		return fmt.Errorf("official_go_client_required")
	}
	return RejectUnofficialClient(cliPath)
}
