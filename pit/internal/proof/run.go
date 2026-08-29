package proof

import (
	"context"
	"os/exec"
)

// runWithContext runs the official client and kills it when the deadline passes.
// The combined output is always returned, including on failure, because the
// parse helpers read the root and transaction out of it.
func runWithContext(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	done := make(chan struct{})
	var out []byte
	var err error
	go func() {
		out, err = cmd.CombinedOutput()
		close(done)
	}()
	select {
	case <-done:
		return out, err
	case <-ctx.Done():
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-done
		return out, ctx.Err()
	}
}
