package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WebSource(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch filepath.Ext(path) {
		case ".ts", ".tsx", ".js", ".jsx", ".html":
		default:
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		low := strings.ToLower(string(b))
		for _, w := range []string{"session_key", "private_key", "hl_secret", "mnemonic"} {
			if strings.Contains(low, w) {
				return fmt.Errorf("forbidden_token %s in %s", w, path)
			}
		}
		return nil
	})
}
