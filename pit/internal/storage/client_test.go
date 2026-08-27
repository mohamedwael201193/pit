package storage

import "testing"

func TestRejectUnofficialClient(t *testing.T) {
	if err := RejectUnofficialClient("upload.ts"); err == nil {
		t.Fatal("ts")
	}
	if err := RejectUnofficialClient("/usr/bin/0g-storage-client"); err != nil {
		t.Fatal(err)
	}
}

func TestDownloadMustProve(t *testing.T) {
	args, err := DownloadArgs(Job{
		CLI: "/usr/bin/0g-storage-client", RPC: "https://evmrpc.0g.ai", Indexer: "https://indexer-storage-turbo.0g.ai",
		KeyHex: "0x" + "aa" + repeat("bb", 31), Root: "0xabc123def0", OutPath: "/tmp/o",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !contains(args, "--root") || !contains(args, "--file") || contains(args, "--rpc") {
		t.Fatal(args)
	}
	if err := DownloadMustProve(args); err != nil {
		t.Fatal(err)
	}
	if err := DownloadMustProve([]string{"download"}); err == nil {
		t.Fatal("proof")
	}
}

func TestRejectBadRoot(t *testing.T) {
	if err := RejectBadRoot("abc"); err == nil {
		t.Fatal("root")
	}
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
