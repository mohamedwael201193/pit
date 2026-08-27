package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0gfoundation/0g-pc-e2ee/client/sig"
	"github.com/0gfoundation/0g-pc-e2ee/protocol/crypto"
	"github.com/0gfoundation/0g-pc-e2ee/protocol/proof"
	"github.com/0gfoundation/0g-pc-e2ee/protocol/wire"
)

func main() {
	authPath := flag.String("auth", "", "Direct auth json (never a Router key)")
	promptPath := flag.String("prompt", "", "utf-8 prompt (hashed in evidence, not copied)")
	outPath := flag.String("out", "", "evidence json (no plaintext book)")
	role := flag.String("role", "", "researcher|challenger|risk")
	tamperCT := flag.Bool("tamper-ciphertext", false, "mutate response ciphertext then verify (expect FAIL)")
	wrongSigner := flag.Bool("wrong-signer", false, "verify against 0x000...001 (expect FAIL)")
	flag.Parse()
	if *authPath == "" || *promptPath == "" || *outPath == "" {
		fmt.Println("usage: pit-sealer -auth -prompt -out -role researcher|challenger|risk")
		os.Exit(2)
	}
	if err := requireRole(*role); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	prompt, err := os.ReadFile(*promptPath)
	if err != nil {
		fmt.Println("READ_PROMPT", err)
		os.Exit(1)
	}
	a, err := loadAuth(*authPath)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "ROUTER_DOWNGRADE_DENIED" {
			os.Exit(10)
		}
		if err.Error() == "NOT_TEEML" {
			os.Exit(11)
		}
		os.Exit(1)
	}
	code := runSeal(a, prompt, *role, *outPath, *tamperCT, *wrongSigner)
	os.Exit(code)
}

func runSeal(a authFile, prompt []byte, role, outPath string, tamperCT, wrongSigner bool) int {
	base := strings.TrimRight(a.URL, "/")
	t0 := time.Now()
	result := map[string]any{
		"role":          role,
		"provider":      a.Provider,
		"url":           base,
		"model":         a.Model,
		"teeSigner":     a.TeeSigner,
		"verifiability": a.Verifiability,
		"prompt_sha256": sha256Hex(prompt),
		"prompt_bytes":  len(prompt),
		"started_unix":  t0.Unix(),
	}
	st, body, _, err := getJSON(base+"/v1/e2ee/pubkey", nil)
	if err != nil {
		fmt.Println("PUBKEY", err)
		_ = writeEvidence(outPath, result)
		return 1
	}
	result["pubkey_http"] = st
	var pk struct {
		EncPub        string `json:"enc_pub"`
		KeyID         string `json:"key_id"`
		SignerAddress string `json:"signer_address"`
	}
	if err := json.Unmarshal(body, &pk); err != nil {
		fmt.Println("PUBKEY_JSON", err)
		_ = writeEvidence(outPath, result)
		return 1
	}
	result["key_id"] = pk.KeyID
	result["pubkey_signer"] = pk.SignerAddress
	if err := requirePubKeyMatchesOnchain(pk.SignerAddress, a.TeeSigner); err != nil {
		result["verify_e2ee"] = "FAIL"
		result["verify_err"] = err.Error()
		_ = writeEvidence(outPath, result)
		fmt.Println("VERIFY_FAIL", err)
		return 6
	}
	encPub, err := base64.RawURLEncoding.DecodeString(pk.EncPub)
	if err != nil {
		encPub, err = base64.StdEncoding.DecodeString(pk.EncPub)
	}
	if err != nil {
		fmt.Println("ENC_PUB", err)
		return 1
	}
	ephPriv, ephPub, err := crypto.GenerateRecipientKey()
	if err != nil {
		fmt.Println("EPH", err)
		return 1
	}
	modelJSON, _ := json.Marshal(a.Model)
	msgJSON, _ := json.Marshal([]map[string]string{{"role": "user", "content": string(prompt)}})
	req := wire.Request{
		"model":    json.RawMessage(modelJSON),
		"messages": json.RawMessage(msgJSON),
	}
	signerHint := pk.SignerAddress
	if signerHint == "" {
		signerHint = a.TeeSigner
	}
	sealed, err := wire.SealRequest(crypto.PublicKey(encPub), req, []string{"messages"}, signerHint, ephPub)
	if err != nil {
		fmt.Println("SEAL", err)
		return 1
	}
	raw, _ := json.Marshal(sealed)
	result["sealed_len"] = len(raw)
	result["sealed_sha256"] = sha256Hex(raw)
	hdr := map[string]string{"Authorization": a.Authorization, "Allow-Fallbacks": "false"}
	st, resp, rh, err := postJSON(base+"/v1/proxy/chat/completions", raw, hdr)
	if err != nil {
		fmt.Println("POST", err)
		return 1
	}
	result["post_http"] = st
	result["response_sha256"] = sha256Hex(resp)
	result["elapsed_ms"] = time.Since(t0).Milliseconds()
	result["finished_unix"] = time.Now().Unix()
	chatID := rh.Get("ZG-Res-Key")
	result["zg_res_key"] = chatID
	if st < 200 || st >= 300 {
		clip := string(resp)
		if len(clip) > 400 {
			clip = clip[:400]
		}
		result["post_err_clip"] = clip
		_ = writeEvidence(outPath, result)
		fmt.Println("POST_FAIL", st)
		return 3
	}
	reqEnv, err := asEnv(raw)
	if err != nil {
		fmt.Println("REQ_ENV", err)
		return 1
	}
	respEnv, err := asEnv(resp)
	if err != nil {
		fmt.Println("RESP_ENV", err)
		return 1
	}
	if chatID == "" {
		var probe map[string]any
		_ = json.Unmarshal(resp, &probe)
		if id, ok := probe["id"].(string); ok {
			chatID = strings.TrimPrefix(id, "chatcmpl-")
			result["zg_res_key_from_id"] = chatID
		}
	}
	if chatID == "" {
		_ = writeEvidence(outPath, result)
		fmt.Println("NO_CHAT_ID")
		return 4
	}
	sigURL := fmt.Sprintf("%s/v1/proxy/signature/%s?model=%s", base, chatID, a.Model)
	st, sigBody, _, err := getJSON(sigURL, hdr)
	result["sig_http"] = st
	if err != nil || st != 200 {
		_ = writeEvidence(outPath, result)
		fmt.Println("SIG_FAIL", st, err)
		return 5
	}
	var cs proof.ChatSignature
	if err := json.Unmarshal(sigBody, &cs); err != nil {
		fmt.Println("SIG_JSON", err)
		return 1
	}
	result["sig_text"] = cs.Text
	result["sig_algo"] = cs.SigningAlgo
	result["sig_self_reported"] = cs.SigningAddress
	result["scheme_ok"] = strings.HasPrefix(cs.Text, "zg-sig-v1/e2ee-ct")

	verifySigner := a.TeeSigner
	if wrongSigner {
		verifySigner = "0x0000000000000000000000000000000000000001"
	}
	verifyEnv := respEnv
	if tamperCT {
		mut := make(map[string]json.RawMessage, len(respEnv))
		for k, v := range respEnv {
			mut[k] = v
		}
		if e2, ok := mut["_e2ee"]; ok {
			var obj map[string]any
			_ = json.Unmarshal(e2, &obj)
			if ct, ok := obj["ciphertext"].(string); ok && len(ct) > 8 {
				obj["ciphertext"] = "AAAA" + ct[4:]
				nb, _ := json.Marshal(obj)
				mut["_e2ee"] = nb
			}
		}
		verifyEnv = mut
		result["tamper_ciphertext"] = true
	}

	err = cs.VerifyE2EE(reqEnv, verifyEnv, verifySigner, sig.Recover)
	if err != nil {
		result["verify_e2ee"] = "FAIL"
		result["verify_err"] = err.Error()
		_ = writeEvidence(outPath, result)
		fmt.Println("VERIFY_FAIL", err)
		if tamperCT || wrongSigner {
			return 0
		}
		return 6
	}
	result["verify_e2ee"] = "OK"
	result["pubkey_matches_onchain_tee"] = strings.EqualFold(pk.SignerAddress, a.TeeSigner)

	opened, err := wire.OpenResponse(ephPriv, wire.Response(respEnv))
	if err != nil {
		result["open_err"] = err.Error()
		_ = writeEvidence(outPath, result)
		fmt.Println("OPEN_FAIL", err)
		return 7
	}
	text := extractText(opened)
	result["output_sha256"] = sha256Hex([]byte(text))
	result["output_chars"] = len(text)
	result["sanitized_output"] = text
	if err := writeEvidence(outPath, result); err != nil {
		fmt.Println("EVIDENCE", err)
		return 1
	}
	fmt.Println("COMMITTEE_OK", role, cs.Text)
	return 0
}
