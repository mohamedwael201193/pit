package compute

import "fmt"

// Progress maps to real backend steps. UI must not animate these on a timer.
type Progress string

const (
	ProgBook    Progress = "PRIVATE_BOOK"
	ProgSeal    Progress = "SEALING"
	ProgTEE     Progress = "TEE"
	ProgSig     Progress = "TEE_SIGNATURE"
	ProgSigner  Progress = "ONCHAIN_SIGNER"
	ProgStorage Progress = "STORAGE"
	ProgReceipt Progress = "RECEIPT"
	ProgCalib   Progress = "CALIBRATION"
)

type Event struct {
	Progress Progress
	OK       bool
	Detail   string
}

type Role string

const (
	Researcher Role = "researcher"
	Challenger Role = "challenger"
	Risk       Role = "risk"
)

func Next(cur Progress) (Progress, error) {
	order := []Progress{ProgBook, ProgSeal, ProgTEE, ProgSig, ProgSigner, ProgStorage, ProgReceipt, ProgCalib}
	for i, p := range order {
		if p == cur && i+1 < len(order) {
			return order[i+1], nil
		}
	}
	return "", fmt.Errorf("no_next")
}

func Envelope(role Role, publicMarket, privateBook []byte) ([]byte, error) {
	switch role {
	case Researcher, Challenger, Risk:
	default:
		return nil, fmt.Errorf("bad_role")
	}
	if len(publicMarket) == 0 || len(privateBook) == 0 {
		return nil, fmt.Errorf("empty_envelope")
	}
	out := make([]byte, 0, 16+len(publicMarket)+len(privateBook))
	out = append(out, []byte(string(role)+"|")...)
	out = append(out, publicMarket...)
	out = append(out, '|')
	out = append(out, privateBook...)
	return out, nil
}
