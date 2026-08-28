package storage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func clientNames() []string {
	return []string{"0g-storage-client.exe", "0g-storage-client"}
}

func discoverClient(roots ...string) string {
	for _, root := range roots {
		if root == "" {
			continue
		}
		for _, n := range clientNames() {
			p := filepath.Join(root, n)
			st, err := os.Stat(p)
			if err == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

func LookCLI() string {
	if env := strings.TrimSpace(os.Getenv("PIT_STORAGE_CLI")); env != "" {
		return env
	}
	var roots []string
	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(exe))
	}
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, wd)
	}
	if la := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); la != "" {
		roots = append(roots, filepath.Join(la, "PIT"))
	}
	if found := discoverClient(roots...); found != "" {
		return found
	}
	if p, err := exec.LookPath("0g-storage-client"); err == nil {
		return p
	}
	if p, err := exec.LookPath("0g-storage-client.exe"); err == nil {
		return p
	}
	return ""
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
