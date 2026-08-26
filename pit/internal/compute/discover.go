package compute

import (
	"os"
	"path/filepath"
	"strings"
)

func DiscoverSealer(roots ...string) string {
	names := []string{"pit-sealer", "pit-sealer.exe"}
	for _, root := range roots {
		if root == "" {
			continue
		}
		for _, n := range names {
			p := filepath.Join(root, n)
			st, err := os.Stat(p)
			if err == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

func LookBin() string {
	if env := strings.TrimSpace(os.Getenv("PIT_COMMITTEE_BIN")); env != "" {
		return env
	}
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return DiscoverSealer(
		filepath.Join(wd, "sealer"),
		filepath.Join(wd, "..", "sealer"),
		wd,
	)
}
