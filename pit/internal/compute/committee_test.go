package compute

import "testing"

func TestEnvelopeAndProgress(t *testing.T) {
	if _, err := Envelope("nope", []byte("{}"), []byte("{}")); err == nil {
		t.Fatal("role")
	}
	b, err := Envelope(Researcher, []byte(`{"mark":1}`), []byte(`{"clip":10}`))
	if err != nil || len(b) == 0 {
		t.Fatal(err)
	}
	n, err := Next(ProgBook)
	if err != nil || n != ProgSeal {
		t.Fatal(n, err)
	}
	if _, err := Next(ProgCalib); err == nil {
		t.Fatal("end")
	}
}

func TestVerifyFailClosed(t *testing.T) {
	if err := RequireScheme("zg-sig-v1/plain"); err == nil {
		t.Fatal("plain")
	}
	if err := RequireScheme("zg-sig-v1/e2ee-ct:aa:bb"); err != nil {
		t.Fatal(err)
	}
	if err := RequireSigner("0xabc", "0xdef"); err == nil {
		t.Fatal("signer")
	}
	if err := RequireSigner("0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9", "0xa46ea4fc5889ad35a1487e1ed04dccfa872146b9"); err != nil {
		t.Fatal(err)
	}
	if err := RejectPlaintextFallback("true"); err == nil {
		t.Fatal("fallback")
	}
	if err := DenyRouter("https://router-api.0g.ai/v1"); err == nil {
		t.Fatal("router")
	}
}
