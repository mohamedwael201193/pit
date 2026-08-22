package sdk

import "github.com/mohamedwael201193/pit/internal/engine"

func (c Client) ForecastFor(workspace string, rec engine.StoredForecast) (engine.Forecast, error) {
	return engine.Lookup(workspace, rec)
}
