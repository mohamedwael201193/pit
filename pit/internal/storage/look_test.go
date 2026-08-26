package storage

import (
	"strings"
	"testing"
)

func TestProofJobRequiresOfficialClientAndProof(t *testing.T) {
	_, err := ProofJob("", "https://evmrpc.0g.ai", "https://indexer-storage-turbo.0g.ai", "0x"+strings.Repeat("ab", 32), "0x"+strings.Repeat("cd", 32), "out.bin")
	if err == nil {
		t.Fatal("cli")
	}
	j, err := ProofJob("/usr/bin/0g-storage-client", "https://evmrpc.0g.ai", "https://indexer-storage-turbo.0g.ai", "0x"+strings.Repeat("ab", 32), "0x"+strings.Repeat("cd", 32), "out.bin")
	if err != nil {
		t.Fatal(err)
	}
	args, err := DownloadArgs(j)
	if err != nil {
		t.Fatal(err)
	}
	red := RedactArgs(args)
	joined := strings.Join(red, " ")
	if strings.Contains(joined, strings.Repeat("ab", 32)) {
		t.Fatal("key leaked")
	}
	if !strings.Contains(joined, "--proof") {
		t.Fatal("proof")
	}
}
