package sdk

import "github.com/mohamedwael201193/pit/internal/exec"

func (c Client) PreviewFields() []string {
	return append([]string{}, exec.PreviewBindFields...)
}
