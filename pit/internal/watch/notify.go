package watch

import "fmt"

// Watch never holds an execution session and never posts orders.
func MayPlaceOrder(fromWatch bool) error {
	if fromWatch {
		return fmt.Errorf("watch_cannot_place_orders")
	}
	return nil
}

func NotifyCopy(n int) string {
	return EmptyCopy(n)
}
