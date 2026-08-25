package sdk

import "github.com/mohamedwael201193/pit/internal/engine"

func (c Client) DropModelSize(model map[string]any) error {
	return engine.IgnoreModelSize(model)
}
