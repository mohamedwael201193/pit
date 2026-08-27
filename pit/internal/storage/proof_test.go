package storage

import (
	"strings"
	"testing"
)

func TestProofArgs(t *testing.T) {
	j := Job{CLI: "0g-storage-client", RPC: "https://evmrpc.0g.ai", Indexer: "https://indexer-storage-turbo.0g.ai", KeyHex: "0x" + strings.Repeat("aa", 32), PayerKey: "0x" + strings.Repeat("11", 32), Root: "0x" + strings.Repeat("bb", 32), InputPath: "in.bin", OutPath: "out.bin"}
	up, err := UploadArgs(j)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(up, "--encryption-key") || !contains(up, "--url") || !contains(up, "--file") || contains(up, "--rpc") {
		t.Fatal(up)
	}
	if contains(RedactArgs(up), strings.Repeat("11", 32)) {
		t.Fatal("payer leaked")
	}
	down, err := DownloadArgs(j)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(down, "--proof") || !contains(down, "--root") || !contains(down, "--file") {
		t.Fatal(down)
	}
	if err := DownloadMustProve(down); err != nil {
		t.Fatal(err)
	}
	j.CLI = "sdk.ts"
	if _, err := UploadArgs(j); err == nil {
		t.Fatal("ts forbidden")
	}
	j.CLI = "0g-storage-client"
	j.PayerKey = j.KeyHex
	if _, err := UploadArgs(j); err == nil {
		t.Fatal("payer must not equal memory key")
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
