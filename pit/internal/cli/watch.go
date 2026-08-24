package cli

import "github.com/mohamedwael201193/pit/internal/watch"

func WatchCopy(n int) string {
	return watch.NotifyCopy(n)
}

func WatchMayPlace() error {
	return watch.MayPlaceOrder(true)
}
