package hl

import "fmt"

func IndexInUniverse(universe []map[string]any, coin string) (int, error) {
	if coin == "" {
		return -1, fmt.Errorf("unknown coin")
	}
	for i, u := range universe {
		name, _ := u["name"].(string)
		if name == coin {
			return i, nil
		}
	}
	return -1, fmt.Errorf("unknown coin")
}
