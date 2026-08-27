package exec

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func testSignedOrder(t *testing.T) hl.Envelope {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := hl.BuildOrder(1, true, "2500", "0.004", "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	env, err := hl.SignL1(key, raw, 1700000000000, false)
	if err != nil {
		t.Fatal(err)
	}
	return env
}

func TestPostSignedRefusesUnsignedAndUnlinked(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	t.Cleanup(srv.Close)
	e := NewExchange(srv.URL)
	env := testSignedOrder(t)
	if _, err := PostSigned(e, hl.Envelope{Action: env.Action}, true, "h", "h"); err == nil {
		t.Fatal("unsigned")
	}
	if _, err := PostSigned(e, env, false, "h", "h"); err == nil {
		t.Fatal("unlinked")
	}
	if _, err := PostSigned(e, env, true, "aaa", "bbb"); err == nil {
		t.Fatal("hash")
	}
	if hits.Load() != 0 {
		t.Fatal("posted")
	}
}

func TestPostSignedPostsEnvelopeWhenLinked(t *testing.T) {
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok","response":{"type":"order","data":{"statuses":[{"resting":{"oid":42}}]}}}`))
	}))
	t.Cleanup(srv.Close)
	e := NewExchange(srv.URL)
	env := testSignedOrder(t)
	got, err := PostSigned(e, env, true, "bound", "bound")
	if err != nil {
		t.Fatal(err)
	}
	if ReceiptOID(got) != "42" {
		t.Fatalf("oid %s", ReceiptOID(got))
	}
	var posted hl.Envelope
	if err := json.Unmarshal(body, &posted); err != nil || !posted.Signed() {
		t.Fatalf("envelope %s %v", body, err)
	}
	if err := e.Guard(posted.Action); err != nil {
		t.Fatal(err)
	}
}

func TestReceiptOIDEmptyOnError(t *testing.T) {
	if ReceiptOID([]byte(`{"status":"err","response":"bad"}`)) != "" {
		t.Fatal("err")
	}
}
