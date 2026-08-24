package sdk

import "github.com/mohamedwael201193/pit/internal/calib"

func (c Client) HealthCard(n int) calib.Health {
	return calib.Card(nil, n)
}
