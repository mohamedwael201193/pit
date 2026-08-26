package redteam

func BookNeverInHealth(body map[string]any) bool {
	if body == nil {
		return true
	}
	_, book := body["private_book"]
	_, strat := body["strategy"]
	trade, _ := body["trade"].(bool)
	sign, _ := body["sign"].(bool)
	return !book && !strat && !trade && !sign
}
