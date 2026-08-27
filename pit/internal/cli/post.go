package cli

import (
	"github.com/mohamedwael201193/pit/internal/config"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func PostLinked(network string, env hl.Envelope, linked bool, previewHash string) ([]byte, error) {
	net, err := config.ParseNetwork(network)
	if err != nil {
		return nil, err
	}
	return pitexec.PostSigned(pitexec.NewExchange(config.For(net).HLExchange), env, linked, previewHash, previewHash)
}
