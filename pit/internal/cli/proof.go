package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mohamedwael201193/pit/internal/storage"
)

type ProofFlags struct {
	Root    string
	Out     string
	KeyFile string
}

func ParseProofFlags(args []string) (ProofFlags, error) {
	var f ProofFlags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--root":
			i++
			if i < len(args) {
				f.Root = args[i]
			}
		case "--out":
			i++
			if i < len(args) {
				f.Out = args[i]
			}
		case "--key-file":
			i++
			if i < len(args) {
				f.KeyFile = args[i]
			}
		}
	}
	if f.Root == "" || f.Out == "" || f.KeyFile == "" {
		return ProofFlags{}, fmt.Errorf("proof download needs --root, --out, and a workspace key file. PIT does not use a global memory key.")
	}
	return f, nil
}

func LoadProofKey(path, envKey string) (string, error) {
	if strings.TrimSpace(path) == "" {
		if strings.TrimSpace(envKey) != "" {
			return "", storage.RefuseGlobalMemoryKey(envKey)
		}
		return "", fmt.Errorf("key_file_required")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	key := strings.TrimSpace(string(b))
	if err := storage.RequireHexKey(key); err != nil {
		return "", err
	}
	return key, nil
}
