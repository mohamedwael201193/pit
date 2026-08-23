package obs

import "testing"

func TestRefuseSensitiveLog(t *testing.T) {
	if err := RefuseSensitiveLog("phase=SEALING"); err != nil {
		t.Fatal(err)
	}
	if err := RefuseSensitiveLog("session_key leaked"); err == nil {
		t.Fatal("key")
	}
	if err := RefuseSensitiveLog("private book dump"); err == nil {
		t.Fatal("book")
	}
}
