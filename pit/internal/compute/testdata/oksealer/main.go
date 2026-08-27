package main

import (
	"os"
)

func main() {
	out := ""
	for i, a := range os.Args {
		if a == "-out" && i+1 < len(os.Args) {
			out = os.Args[i+1]
		}
	}
	body := `{"verify_e2ee":"OK","sig_text":"zg-sig-v1/e2ee-ct:aa","pubkey_signer":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9","teeSigner":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"}`
	if out != "" {
		_ = os.WriteFile(out, []byte(body), 0o600)
	}
	_, _ = os.Stdout.WriteString("COMMITTEE_OK researcher zg-sig-v1/e2ee-ct:aa\n")
}
