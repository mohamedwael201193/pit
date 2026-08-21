package ui

var you = []string{
	"YOUR WALLET",
	"YOUR TRADING ACCOUNT",
	"YOUR SESSION",
	"YOUR POLICY",
	"YOUR MONEY",
	"YOUR DESK",
	"YOUR NETWORK",
}

func Labels() []string {
	out := make([]string, len(you))
	copy(out, you)
	return out
}

func HasSeedField(form string) bool {
	return false
}
