package phase

import "fmt"

var Order = []string{
	Connecting, Authenticating, Sealing, Researching, Challenging, AssessingRisk,
	Scoring, PolicyCheck, WaitingForUser, Signing, Submitting, Confirming,
	Executed, Verifying, Resolved, Calibrated,
}

func Next(cur string) (string, error) {
	if cur == Failed {
		return "", fmt.Errorf("terminal")
	}
	for i, s := range Order {
		if s == cur && i+1 < len(Order) {
			return Order[i+1], nil
		}
	}
	if cur == Calibrated {
		return "", fmt.Errorf("terminal")
	}
	return "", fmt.Errorf("unknown_phase")
}
