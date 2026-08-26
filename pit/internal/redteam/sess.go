package redteam

import (
	"strings"

	"github.com/mohamedwael201193/pit/internal/cli"
)

func SessionMetaLeak(raw string) error {
	return cli.RefuseSessionSecret([]byte(raw))
}

func SessionMetaLooksPublic(raw string) bool {
	return !strings.Contains(strings.ToLower(raw), "private") && SessionMetaLeak(raw) == nil
}
