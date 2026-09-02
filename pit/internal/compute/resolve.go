package compute

import (
	"fmt"
	"strings"
)

const preferredPrivateModel = "glm-5.3"

func teemlSet(models []ProviderModel) map[string]struct{} {
	out := map[string]struct{}{}
	for _, m := range models {
		if RequireTeeML(m.Verifiability) != nil {
			continue
		}
		id := strings.ToLower(strings.TrimSpace(m.ID))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

// PickDirectModel chooses the sealed-book SKU from on-chain getService plus this
// provider's /v1/models. glm-5.3 is used only when THIS Direct provider lists it
// as TeeML. Router TeeTLS glm-5.3 is never selected. If glm-5.3 TeeML is present,
// glm-5.2 is not used.
func PickDirectModel(live LiveService, models []ProviderModel) (string, error) {
	if err := RequireTeeML(live.Verifiability); err != nil {
		return "", err
	}
	teeml := teemlSet(models)
	if _, ok := teeml[preferredPrivateModel]; ok {
		return preferredPrivateModel, nil
	}
	liveModel := strings.TrimSpace(live.Model)
	if liveModel == "" {
		return "", fmt.Errorf("private_teeml_sku_unavailable")
	}
	if RequireTeeML(live.Verifiability) != nil {
		return "", fmt.Errorf("NOT_TEEML")
	}
	if len(models) > 0 {
		if _, ok := teeml[strings.ToLower(liveModel)]; !ok {
			return "", fmt.Errorf("private_teeml_sku_unavailable")
		}
	}
	return liveModel, nil
}
