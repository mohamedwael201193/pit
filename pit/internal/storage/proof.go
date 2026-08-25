package storage

import (
	"fmt"
	"os/exec"
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
	if err := RejectUnofficialClient(j.CLI); err != nil {
		return nil, err
	}
	if err := RequireHexKey(j.KeyHex); err != nil {
		return nil, err
	}
	if j.InputPath == "" || j.RPC == "" || j.Indexer == "" {
		return nil, fmt.Errorf("incomplete_upload")
	}
	args := []string{
		"upload", "--rpc", j.RPC, "--indexer", j.Indexer,
		"--encryption-key", j.KeyHex, j.InputPath,
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
	if j.OutPath == "" {
		return nil, fmt.Errorf("incomplete_download")
	}
	args := []string{
		"download", "--rpc", j.RPC, "--indexer", j.Indexer,
		"--encryption-key", j.KeyHex, "--proof", j.Root, j.OutPath,
	}
	if err := DownloadMustProve(args); err != nil {
		return nil, err
	}
	return args, nil
}

func Command(j Job, args []string) *exec.Cmd {
	return exec.Command(j.CLI, args...)
}
