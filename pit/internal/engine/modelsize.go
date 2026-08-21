package engine

import (
	"encoding/json"
	"fmt"
	"math"
)

func IgnoreModelSize(model map[string]any) error {
	if model == nil {
		return nil
	}
	for _, k := range []string{"sizeUsd", "sz", "notional", "qty"} {
		if _, ok := model[k]; ok {
			delete(model, k)
		}
	}
	return nil
}

func RejectModelSize(raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	v, ok := m["sizeUsd"]
	if !ok {
		return nil
	}
	switch t := v.(type) {
	case float64:
		if math.IsNaN(t) || math.IsInf(t, 0) {
			return fmt.Errorf("model_size_ignored_invalid")
		}
	}
	return nil
}
