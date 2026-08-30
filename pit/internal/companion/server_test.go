package companion

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func local(r *http.Request) *http.Request {
	r.RemoteAddr = "127.0.0.1:9"
	return r
}

func TestPairingIssuesDeviceNotSession(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["sign"] != false || got["canSign"] != false {
		t.Fatalf("%+v", got)
	}
	tok, _ := got["device"].(string)
	if len(tok) < 32 {
		t.Fatal("token")
	}
	if strings.Contains(rec.Body.String(), "private") {
		t.Fatal("secret")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("replay")
	}
}

func TestMaliciousOriginCannotPair(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal(rec.Code)
	}
}

func TestAuthorizeExportDenied(t *testing.T) {
	h := New(t.TempDir())
	for _, path := range []string{"/authorize", "/export", "/session"} {
		req := local(httptest.NewRequest(http.MethodPost, path, nil))
		req.Header.Set("Origin", "https://pit0g.vercel.app")
		rec := httptest.NewRecorder()
		h.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("%s %d", path, rec.Code)
		}
	}
}

func TestListenRejectsPublicBind(t *testing.T) {
	t.Setenv("PIT_COMPANION_ADDR", "0.0.0.0:17373")
	if _, err := ListenAddr(); err == nil {
		t.Fatal("public bind")
	}
}

func TestWrongCodeDenied(t *testing.T) {
	h := New(t.TempDir())
	body, _ := json.Marshal(map[string]string{"code": "AAAAAAAA"})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal(rec.Code)
	}
}

func TestLocalStatusAllowsLoopbackOmitsCode(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	req := local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), code) {
		t.Fatal("code leak")
	}
	req = local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	req.Header.Set("Origin", "https://evil.example")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal(rec.Code)
	}
}

func TestHealthOmitsPairingCode(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	req := local(httptest.NewRequest(http.MethodGet, "/health", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
	if strings.Contains(rec.Body.String(), code) {
		t.Fatal("code leak")
	}
}

func TestLocalCodeDesktopOnly(t *testing.T) {
	h := New(t.TempDir())
	want, _ := h.Code()

	req := local(httptest.NewRequest(http.MethodGet, "/local/code", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not read pairing code")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/code", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), want) {
		t.Fatal("desktop code")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/code", nil))
	req.Header.Set("Origin", "http://tauri.localhost")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal("windows webview pairing")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/doctor", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web doctor")
	}
}

func TestLocalCodeRotateAndPairingState(t *testing.T) {
	h := New(t.TempDir())
	first, _ := h.Code()
	req := local(httptest.NewRequest(http.MethodPost, "/local/code/rotate", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["sign"] == true || got["trade"] == true {
		t.Fatal(got)
	}
	next, _ := got["code"].(string)
	if next == "" || next == first {
		t.Fatalf("rotate %q -> %q", first, next)
	}
	req = local(httptest.NewRequest(http.MethodPost, "/local/code/rotate", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not rotate pairing code")
	}
}

func TestCatalogListingCannotBecomeChatModel(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/models", strings.NewReader(`{"model":"0gm-1.0-35b-a3b"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "catalog_listing_not_inference") {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"picked":"host-parsed"`) {
		t.Fatal(rec.Body.String())
	}
}

func TestLocalStatusVersionNoSecret(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["sign"] == true || got["trade"] == true {
		t.Fatal(got)
	}
	if got["version"] != "0.9.0" {
		t.Fatalf("version %v", got["version"])
	}
}

func TestDesktopInitSessionPolicyUserB(t *testing.T) {
	h := New(t.TempDir())
	a := "0x1111111111111111111111111111111111111111"
	b := "0x2222222222222222222222222222222222222222"

	req := local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"`+a+`","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "private") {
		t.Fatal("secret")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"`+b+`","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("user B %d %s", rec.Code, rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"`+a+`","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not init")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/session", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var sess map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &sess); err != nil {
		t.Fatal(err)
	}
	if sess["withdraw"] != false || sess["sign"] != false || sess["agent"] == nil || sess["agent"] == "" {
		t.Fatal(sess)
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "private_key") {
		t.Fatal("key leak")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/session", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web session")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/policy", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), a) {
		t.Fatal("wallet missing")
	}
}

func TestPolicyPinClampAndWebDenied(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/policy", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	body := `{"maxClipUsd":12,"dailyLossUsd":50,"maxLeverage":20,"allowedAssets":["ETH","BTC"],"maxOpenPositions":2,"maxSlippageBps":80,"cooldownSeconds":0,"sessionTtlSeconds":3600,"maxUncertainty":1,"minLiquidityUsd":0,"maxConsecutiveLosses":3}`
	req = local(httptest.NewRequest(http.MethodPost, "/local/policy", strings.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"maxLeverage":20`) {
		t.Fatal("leverage must stay 1")
	}
	if !strings.Contains(rec.Body.String(), `"maxOpenPositions":2`) || !strings.Contains(rec.Body.String(), `"maxClipUsd":12`) {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/policy", strings.NewReader(`{"maxClipUsd":1000,"maxLeverage":50}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("web pin %d", rec.Code)
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/policy", strings.NewReader(`{"maxClipUsd":1000,"maxLeverage":50,"dailyLossUsd":50,"allowedAssets":["ETH"],"maxOpenPositions":1,"maxSlippageBps":80,"sessionTtlSeconds":3600,"maxUncertainty":1,"maxConsecutiveLosses":3}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"maxClipUsd":1000`) || strings.Contains(rec.Body.String(), `"maxLeverage":50`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"maxClipUsd":50`) {
		t.Fatal(rec.Body.String())
	}
}

func TestBrowserBindRequiresDeviceAndIsolatesUserB(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var paired map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &paired); err != nil {
		t.Fatal(err)
	}
	tok, _ := paired["device"].(string)

	req = local(httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatal("no bearer")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"wallet":"0x2222222222222222222222222222222222222222","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("user B bind %d %s", rec.Code, rec.Body.String())
	}
}

func TestDirectIntentNeverLeaksToken(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/direct-intent", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	body := rec.Body.String()
	if strings.Contains(strings.ToLower(body), "app-sk-") {
		t.Fatal("token leak")
	}
	if !strings.Contains(body, `"digest"`) || !strings.Contains(body, `"message"`) {
		t.Fatal(body)
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/direct-intent", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not read local direct intent")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/direct/intent", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatal("device required")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/research", strings.NewReader(`{"coin":"ETH"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web research")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/research/start", strings.NewReader(`{"coin":"ETH"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web research start")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web research status")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/authorize", strings.NewReader(`{"typed":"AUTHORIZE"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web authorize")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/cancel", strings.NewReader(`{"typed":"AUTHORIZE"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web cancel")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("status leak")
	}
}

func TestResearchStartReturnsImmediately(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/research/start", strings.NewReader(`{"coin":"ETH"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	start := time.Now()
	h.Handler().ServeHTTP(rec, req)
	if time.Since(start) > 2*time.Second {
		t.Fatal("start blocked on sealed ask")
	}
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("start leak")
	}
	if !strings.Contains(rec.Body.String(), `"running"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"job_id"`) {
		t.Fatal("job_id")
	}
	req = local(httptest.NewRequest(http.MethodPost, "/local/research/cancel", strings.NewReader(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	for i := 0; i < 50; i++ {
		h.researchMu.Lock()
		done := !h.job.running
		h.researchMu.Unlock()
		if done {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestResearchStatusOmitsEvidenceWhileRunning(t *testing.T) {
	h := New(t.TempDir())
	h.researchMu.Lock()
	h.job = researchJob{ID: "job-1", running: true, stage: "CONTACTING_PRIVATE_PROVIDER", coin: "ETH", started: time.Now()}
	h.researchMu.Unlock()
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"evidence"`) {
		t.Fatal("status must stay tiny while running")
	}
	if !strings.Contains(rec.Body.String(), `"CONTACTING_PRIVATE_PROVIDER"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("leak")
	}
}

func TestResearchStatusIncludesCompactRolesWhileRunning(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	ev := `{"sign":false,"trade":false,"roles":[{"role":"researcher","verify_e2ee":"OK","pubkey_signer":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9","teeSigner":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9","proposed_side":"buy"}]}`
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), []byte(ev), 0o600); err != nil {
		t.Fatal(err)
	}
	h.researchMu.Lock()
	h.job = researchJob{ID: "job-2", running: true, stage: "RISK", coin: "ETH", started: time.Now()}
	h.researchMu.Unlock()
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"evidence"`) {
		t.Fatal("status must stay tiny while running")
	}
	if !strings.Contains(rec.Body.String(), `"verify_e2ee":"OK"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"RISK"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("leak")
	}
}

func TestClassifyResearch(t *testing.T) {
	if classifyResearch("unbound") != "WORKSPACE_NOT_BOUND" {
		t.Fatal("unbound")
	}
	if classifyResearch("direct_token_required") != "DIRECT_NOT_AUTHORIZED" {
		t.Fatal("token")
	}
	if classifyResearch("policy_changed") != "POLICY_REJECTED" {
		t.Fatal("pin")
	}
	if classifyResearch("TEE_VERIFY_FAIL") != "TEE_SIGNATURE_INVALID" {
		t.Fatal("tee")
	}
	if classifyResearch("empty_envelope") != "HL_MARKET_UNAVAILABLE" {
		t.Fatal("book")
	}
	if classifyResearch("SPONSOR_QUOTA") != "SPONSOR_QUOTA" {
		t.Fatal("quota")
	}
	if classifyResearch("risk_killed") != "risk_killed" {
		t.Fatal("risk")
	}
	if classifyResearch("challenger_killed") != "challenger_killed" {
		t.Fatal("challenger")
	}
	if classifyResearch("committee_incomplete") != "COMMITTEE_INCOMPLETE" {
		t.Fatal("incomplete")
	}
	if classifyResearch("DIRECT_PROVIDER_TIMEOUT") != "DIRECT_PROVIDER_TIMEOUT" {
		t.Fatal("timeout")
	}
	if classifyResearch("direct_signature_http") != "DIRECT_PROVIDER_TIMEOUT" {
		t.Fatal("sig http")
	}
	if classifyResearch("direct_ledger") != "DIRECT_CREDIT_INSUFFICIENT" {
		t.Fatal("ledger")
	}
	if classifyResearch("sealer_runtime") != "DIRECT_PROVIDER_UNAVAILABLE" {
		t.Fatal("runtime")
	}
}

func TestResearchStatusHydratesRolesFromEvidence(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	ev := `{"sign":false,"trade":false,"error":"direct_ledger","roles":[{"role":"researcher","verify_e2ee":"OK","pubkey_signer":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9","teeSigner":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"}]}`
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), []byte(ev), 0o600); err != nil {
		t.Fatal(err)
	}
	h.researchMu.Lock()
	h.job = researchJob{ID: "job-1", done: true, stage: "STOPPED", err: "DIRECT_PROVIDER_UNAVAILABLE", coin: "ETH"}
	h.researchMu.Unlock()
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"verify_e2ee":"OK"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"DIRECT_PROVIDER_UNAVAILABLE"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("leak")
	}
}

func TestResearchStatusKeepsPersistedPreview(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	ev := `{"sign":false,"trade":false,"roles":[{"role":"researcher","verify_e2ee":"OK","proposed_side":"buy"},{"role":"challenger","verify_e2ee":"OK","survives":true,"kill":false},{"role":"risk","verify_e2ee":"OK","survives":true,"kill":false}]}`
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), []byte(ev), 0o600); err != nil {
		t.Fatal(err)
	}
	h.researchMu.Lock()
	h.job = researchJob{
		ID:          "job-keep",
		done:        true,
		stage:       "READY",
		coin:        "ETH",
		eligible:    true,
		previewHash: "0xabcpreview",
		preview:     map[string]any{"eligible": true, "market": "hyperliquid:perp:ETH", "side": "buy", "sz": 0.0041, "limitPx": "2482.2", "hash": "0xabcpreview"},
	}
	h.researchMu.Unlock()
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"verify_e2ee":"OK"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"0xabcpreview"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"2482.2"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"evidence"`) {
		t.Fatal("status must stay tiny when done")
	}
	if !strings.Contains(rec.Body.String(), `"seq"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"heartbeat_unix_ms"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("leak")
	}
}

func TestResearchResultIncludesEvidence(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	ev := `{"sign":false,"trade":false,"roles":[{"role":"researcher","verify_e2ee":"OK","elapsed_ms":8790},{"role":"challenger","verify_e2ee":"OK","survives":true},{"role":"risk","verify_e2ee":"OK","survives":true,"kill":false,"elapsed_ms":5381}]}`
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), []byte(ev), 0o600); err != nil {
		t.Fatal(err)
	}
	h.researchMu.Lock()
	h.job = researchJob{
		ID: "job-result", done: true, stage: "READY", coin: "ETH",
		previewHash: "0xabcpreview",
		preview:     map[string]any{"eligible": true, "hash": "0xabcpreview", "limitPx": "2482.2"},
	}
	h.researchMu.Unlock()
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/result", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"evidence"`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"elapsed_ms":5381`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
		t.Fatal("leak")
	}
}

func TestLoadJobHydratesVerifiedCommittee(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	ev := `{"sign":false,"trade":false,"roles":[{"role":"researcher","verify_e2ee":"OK"},{"role":"challenger","verify_e2ee":"OK"},{"role":"risk","verify_e2ee":"OK"}]}`
	if err := os.WriteFile(filepath.Join(dir, "last-research.json"), []byte(ev), 0o600); err != nil {
		t.Fatal(err)
	}
	job := `{"id":"job-orphan","pid":1,"running":true,"done":false,"stage":"RISK","coin":"ETH","sign":false,"trade":false}`
	if err := os.WriteFile(filepath.Join(dir, "research-job.json"), []byte(job), 0o600); err != nil {
		t.Fatal(err)
	}
	h.researchMu.Lock()
	h.loadJobLocked()
	if h.job.running || !h.job.done || h.job.stage != "READY" || h.job.err != "" {
		t.Fatalf("job %#v", h.job)
	}
	h.researchMu.Unlock()
	req := local(httptest.NewRequest(http.MethodGet, "/local/research/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "COMPANION_NOT_RUNNING") {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"evidence"`) {
		t.Fatal("status must stay tiny")
	}
	if !strings.Contains(rec.Body.String(), `"READY"`) {
		t.Fatal(rec.Body.String())
	}
}

func TestLocalKillDeniedToWebsite(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	req = local(httptest.NewRequest(http.MethodPost, "/local/kill", strings.NewReader(`{"on":true}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web kill")
	}
	req = local(httptest.NewRequest(http.MethodPost, "/local/kill", strings.NewReader(`{"on":true}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"kill":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestLocalIdentitySplitsMintAndTransfer(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodGet, "/local/identity", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"itransfer":"UNAVAILABLE"`) {
		t.Fatal(body)
	}
	if !strings.Contains(body, "Mint is optional") {
		t.Fatal(body)
	}
	if strings.Contains(body, `"itransfer":"LIVE"`) {
		t.Fatal(body)
	}
	req = local(httptest.NewRequest(http.MethodGet, "/local/identity", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web identity")
	}
}
