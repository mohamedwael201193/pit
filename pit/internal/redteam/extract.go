package redteam

import "github.com/mohamedwael201193/pit/internal/scan"

func BrowserExtract(root string) error {
	return scan.WebSource(root)
}
