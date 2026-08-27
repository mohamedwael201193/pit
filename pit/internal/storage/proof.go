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
	Flow      string
	KeyHex    string
	PayerKey  string
	Root      string
	InputPath string
	OutPath   string
}

func UploadArgs(j Job) ([]string, error) {
	if err := RejectUnofficialClient(j.CLI); err != nil {
		return nil, err
	}
	if err := RequireHexKey(j.KeyHex); err != nil {
		return nil, err
	}
	payer, err := NormalizePayerKey(j.PayerKey)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(payer, j.KeyHex) {
		return nil, fmt.Errorf("payer_must_not_be_memory_key")
	}
	if j.InputPath == "" || j.RPC == "" || j.Indexer == "" {
		return nil, fmt.Errorf("incomplete_upload")
	}
	args := []string{
		"upload",
		"--url", j.RPC,
		"--indexer", j.Indexer,
		"--file", j.InputPath,
		"--encryption-key", j.KeyHex,
		"--key", payer,
	}
	if strings.TrimSpace(j.Flow) != "" {
		args = append(args, "--flow-address", j.Flow)
	}
	if err := UploadMustEncrypt(args); err != nil {
		return nil, err
	}
	return args, nil
}

func DownloadArgs(j Job) ([]string, error) {
	if err := RejectUnofficialClient(j.CLI); err != nil {
		return nil, err
	}
	if err := RequireHexKey(j.KeyHex); err != nil {
		return nil, err
	}
	if err := RejectBadRoot(j.Root); err != nil {
		return nil, err
	}
	if j.OutPath == "" || j.Indexer == "" {
		return nil, fmt.Errorf("incomplete_download")
	}
	args := []string{
		"download",
		"--indexer", j.Indexer,
		"--file", j.OutPath,
		"--root", j.Root,
		"--encryption-key", j.KeyHex,
		"--proof",
	}
	if err := DownloadMustProve(args); err != nil {
		return nil, err
	}
	return args, nil
}

func Command(j Job, args []string) *exec.Cmd {
	return exec.Command(j.CLI, args...)
}
