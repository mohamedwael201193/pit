package storage

import (
	"fmt"
	"os/exec"
	"strings"
)

type Job struct {
	CLI       string
	RPC       string
	Indexer   string
	KeyHex    string
	Root      string
	InputPath string
	OutPath   string
}

func UploadArgs(j Job) ([]string, error) {
	if err := RequireHexKey(j.KeyHex); err != nil {
		return nil, err
	}
	if j.CLI == "" || strings.HasSuffix(strings.ToLower(j.CLI), ".ts") {
		return nil, fmt.Errorf("official_go_client_required")
	}
	return []string{
		"upload", "--rpc", j.RPC, "--indexer", j.Indexer,
		"--encryption-key", j.KeyHex, j.InputPath,
	}, nil
}

func DownloadArgs(j Job) ([]string, error) {
	if err := RequireHexKey(j.KeyHex); err != nil {
		return nil, err
	}
	if j.Root == "" {
		return nil, fmt.Errorf("root_required")
	}
	return []string{
		"download", "--rpc", j.RPC, "--indexer", j.Indexer,
		"--encryption-key", j.KeyHex, "--proof", j.Root, j.OutPath,
	}, nil
}

func Command(j Job, args []string) *exec.Cmd {
	return exec.Command(j.CLI, args...)
}
