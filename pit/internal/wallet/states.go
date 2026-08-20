package wallet

import "fmt"

// Named states for connect / SIWE. Never collect a seed phrase.
const (
	SignatureDeclined = "SIGNATURE_DECLINED"
	WrongNetwork      = "WRONG_NETWORK"
	SessionExpired    = "SESSION_EXPIRED"
	HLUnfunded        = "HL_UNFUNDED"
	PolicyBlock       = "POLICY_BLOCK"
	TEEVerifyFail     = "TEE_VERIFY_FAIL"
	SeedForbidden     = "SEED_FORBIDDEN"
)

func RejectSeedField(hasSeedInput bool) error {
	if hasSeedInput {
		return fmt.Errorf(SeedForbidden)
	}
	return nil
}

func MapChain(got, want int64) error {
	if got != want {
		return fmt.Errorf(WrongNetwork)
	}
	return nil
}

func MapSignature(ok bool) error {
	if !ok {
		return fmt.Errorf(SignatureDeclined)
	}
	return nil
}
