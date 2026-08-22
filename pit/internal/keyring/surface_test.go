package keyring

import "testing"

func TestRefuseWebAndRedact(t *testing.T) {
	if err := RefuseWeb("web"); err == nil {
		t.Fatal("web")
	}
	if err := RefuseWeb("desktop"); err != nil {
		t.Fatal(err)
	}
	if Redact([]byte("secret")) != "[redacted]" {
		t.Fatal("redact")
	}
}
