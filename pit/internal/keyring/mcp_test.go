package keyring

import "testing"

func TestRefuseMCP(t *testing.T) {
	if err := RefuseMCP(); err == nil {
		t.Fatal("mcp")
	}
	if err := RefuseWeb("desktop"); err != nil {
		t.Fatal(err)
	}
}
