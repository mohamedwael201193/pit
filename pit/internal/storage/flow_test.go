package storage

import "testing"

func TestRequireFlow(t *testing.T) {
	if err := RequireFlow("0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526"); err != nil {
		t.Fatal(err)
	}
	if err := RequireFlow("not-an-address"); err == nil {
		t.Fatal("flow")
	}
}

func TestUploadMustEncrypt(t *testing.T) {
	key := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := UploadMustEncrypt([]string{"upload", "--encryption-key", key, "file"}); err != nil {
		t.Fatal(err)
	}
	if err := UploadMustEncrypt([]string{"upload", "file"}); err == nil {
		t.Fatal("missing key")
	}
}
