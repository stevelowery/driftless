package filter

import "github.com/stevelowery/driftless/internal/api"

func By(orders []*api.Order, include func(*api.Order) bool) []*api.Order {
	var filtered []*api.Order
	for _, order := range orders {
		if include(order) {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

func ByType(orders []*api.Order, orderType api.OrderType) []*api.Order {
	return By(orders, func(order *api.Order) bool {
		return order.Type() == orderType
	})
}
