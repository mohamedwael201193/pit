package storage

import (
	"strings"
	"testing"
)

func TestProofArgs(t *testing.T) {
	j := Job{CLI: "0g-storage-client", RPC: "https://evmrpc.0g.ai", Indexer: "https://indexer-storage-turbo.0g.ai", KeyHex: "0x" + strings.Repeat("aa", 32), Root: "0x" + strings.Repeat("bb", 32), InputPath: "in.bin", OutPath: "out.bin"}
	up, err := UploadArgs(j)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(up, "--encryption-key") {
		t.Fatal(up)
	}
	down, err := DownloadArgs(j)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(down, "--proof") {
		t.Fatal(down)
	}
	if err := DownloadMustProve(down); err != nil {
		t.Fatal(err)
	}
	j.CLI = "sdk.ts"
	if _, err := UploadArgs(j); err == nil {
		t.Fatal("ts forbidden")
	}
}

func contains(ss []string, w string) bool {
	for _, s := range ss {
		if s == w {
			return true
		}
	}
	return false
}
