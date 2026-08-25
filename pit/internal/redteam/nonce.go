package redteam

import "github.com/mohamedwael201193/pit/internal/siwe"

func NonceReplayDenied(used map[string]struct{}, nonce string) bool {
	return siwe.RefuseNonceReplay(used, nonce) != nil
}
