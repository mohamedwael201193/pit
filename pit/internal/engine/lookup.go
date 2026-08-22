package engine

import "fmt"

type StoredForecast struct {
	Workspace string
	ID        string
	Forecast  Forecast
}

func Lookup(ownerWorkspace string, rec StoredForecast) (Forecast, error) {
	if rec.ID == "" || rec.Workspace == "" {
		return Forecast{}, fmt.Errorf("not found")
	}
	if rec.Workspace != ownerWorkspace {
		return Forecast{}, fmt.Errorf("not found")
	}
	return rec.Forecast, nil
}
