package sdk

import "github.com/mohamedwael201193/pit/internal/watch"

func (c Client) WatchCopy(n int) string {
	return watch.Attention(n)
}

func (c Client) WatchCannotTrade() error {
	return watch.MayPlaceOrder(true)
}
