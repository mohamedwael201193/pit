package compute

import (
	"os"
	"path/filepath"
	"strings"
)

func sealerNames() []string {
	return []string{"pit-sealer", "pit-sealer.exe", "pit-sealer-x86_64-pc-windows-msvc.exe", "pit-sealer-x86_64-pc-windows-msvc"}
}

func DiscoverSealer(roots ...string) string {
	for _, root := range roots {
		if root == "" {
			continue
		}
		for _, n := range sealerNames() {
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
	var roots []string
	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(exe))
	}
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots,
			filepath.Join(wd, "sealer"),
			filepath.Join(wd, "..", "sealer"),
			wd,
		)
	}
	return DiscoverSealer(roots...)
}
