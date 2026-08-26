package redteam

import "github.com/mohamedwael201193/pit/internal/exec"

func UnsignedExchange(signed bool) error {
	return exec.RefuseUnsigned(signed)
}
