package companion

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadJobRecoversVerifiedCommittee(t *testing.T) {
	dir := t.TempDir()
	ev := map[string]any{
		"sign":  false,
		"trade": false,
		"roles": []any{
			map[string]any{"role": "researcher", "verify_e2ee": "OK", "pubkey_signer": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9", "teeSigner": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"},
			map[string]any{"role": "challenger", "verify_e2ee": "OK", "pubkey_signer": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9", "teeSigner": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"},
			map[string]any{"role": "risk", "verify_e2ee": "OK", "pubkey_signer": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9", "teeSigner": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"},
		},
	}
	raw, _ := json.Marshal(ev)
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	job, _ := json.Marshal(map[string]any{"id": "job-recover", "running": true, "done": false, "stage": "CONTACTING_PRIVATE_PROVIDER", "coin": "BTC", "pid": 1, "sign": false, "trade": false})
	if err := os.WriteFile(filepath.Join(dir, "research-job.json"), job, 0o600); err != nil {
		t.Fatal(err)
	}
	h := New(dir)
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	body := rec.Body.String()
	if strings.Contains(body, "COMPANION_NOT_RUNNING") {
		t.Fatal(body)
	}
	if !strings.Contains(body, `"READY"`) {
		t.Fatal(body)
	}
	if !strings.Contains(body, `"verify_e2ee":"OK"`) {
		t.Fatal(body)
	}
}

func TestLoadJobMarksFailedWhenNoEvidence(t *testing.T) {
	dir := t.TempDir()
	job, _ := json.Marshal(map[string]any{"id": "job-dead", "running": true, "done": false, "stage": "CONTACTING_PRIVATE_PROVIDER", "coin": "ETH", "pid": 1, "sign": false, "trade": false})
	if err := os.WriteFile(filepath.Join(dir, "research-job.json"), job, 0o600); err != nil {
		t.Fatal(err)
	}
	h := New(dir)
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), "JOB_CRASHED") {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "COMPANION_NOT_RUNNING") {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"FAILED"`) {
		t.Fatal(rec.Body.String())
	}
}

func TestOneRoleNeverVerify(t *testing.T) {
	dir := t.TempDir()
	ev := map[string]any{
		"sign": false, "trade": false,
		"roles": []any{
			map[string]any{"role": "researcher", "verify_e2ee": "OK", "pubkey_signer": "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"},
		},
	}
	raw, _ := json.Marshal(ev)
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	job, _ := json.Marshal(map[string]any{"id": "job-one", "running": false, "done": true, "stage": "READY", "coin": "ETH", "pid": 1, "sign": false, "trade": false})
	if err := os.WriteFile(filepath.Join(dir, "research-job.json"), job, 0o600); err != nil {
		t.Fatal(err)
	}
	h := New(dir)
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	body := rec.Body.String()
	if strings.Contains(body, `"verify":true`) {
		t.Fatal(body)
	}
	if !strings.Contains(body, "COMMITTEE_INCOMPLETE") {
		t.Fatal(body)
	}
	if strings.Contains(body, `"card_title":"RESEARCH COMPLETE"`) {
		t.Fatal(body)
	}
}
