package deskcmd

// Agent is a structured operator payload. Chat/LLM cannot execute from it.
type Agent struct {
	Kind         string      `json:"kind"`
	Executive    string     `json:"executive"`
	Scanned      int         `json:"scanned"`
	Eligible     int         `json:"eligible"`
	Executable   int         `json:"executable"`
	Researched   int         `json:"researched,omitempty"`
	Rejected     int         `json:"rejected,omitempty"`
	Best         string      `json:"best,omitempty"`
	Why          string      `json:"why,omitempty"`
	BuyingPower  float64     `json:"buying_power"`
	PowerSource  string      `json:"power_source,omitempty"`
	AgeMS        int64        `json:"age_ms,omitempty"`
	MinNotional  float64     `json:"min_notional,omitempty"`
	HostNotional float64     `json:"host_notional,omitempty"`
	Mark         float64     `json:"mark,omitempty"`
	Funding      float64     `json:"funding,omitempty"`
	OpenInterest float64     `json:"open_interest,omitempty"`
	Freshness    string      `json:"freshness,omitempty"`
	Coins        []AgentCoin `json:"coins,omitempty"`
}

type AgentCoin struct {
	Coin              string  `json:"coin"`
	Mark              float64 `json:"mark"`
	MinNotional       float64 `json:"min_notional"`
	HostNotional      float64 `json:"host_notional"`
	Funding           float64 `json:"funding"`
	OpenInterest      float64 `json:"open_interest"`
	ExecutionFeasible bool    `json:"execution_feasible"`
	Eligible          bool    `json:"eligible"`
	Why               string  `json:"why,omitempty"`
	ExecWhy           string  `json:"exec_why,omitempty"`
	Block             string  `json:"block,omitempty"`
	Trend             string  `json:"trend,omitempty"`
}
