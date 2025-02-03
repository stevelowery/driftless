package api

type Orders []*Order

type TeamOrders struct {
	Middleton  Orders
	MountHoreb Orders
	Waunakee   Orders
}

func (o Orders) Total() float32 {
	return o.sum(func(order *Order) float32 { return order.Total })
}

func (o Orders) Subtotal() float32 {
	return o.sum(func(order *Order) float32 { return order.Subtotal })
}

func (o Orders) DiscountAmount() float32 {
	return o.sum(func(order *Order) float32 { return order.DiscountAmount })
}

func (o Orders) AmountRefunded() float32 {
	return o.sum(func(order *Order) float32 { return order.AmountRefunded })
}

func (o Orders) Net() float32 {
	return o.sum(func(order *Order) float32 { return order.Net() })
}

func (o Orders) RiderFeeGross() float32 {
	return o.sum(func(order *Order) float32 { return order.RiderFeeGross() })
}

func (o Orders) RiderFeeNet() float32 {
	return o.sum(func(order *Order) float32 { return order.RiderFeeNet() })
}

func (o Orders) BlackhawkFee() float32 {
	return o.sum(func(order *Order) float32 { return order.BlackhawkFee() })
}

func (o Orders) DriftlessDonation() float32 {
	return o.sum(func(order *Order) float32 { return order.DriftlessDonation() })
}

func (o Orders) CORPDonation() float32 {
	return o.sum(func(order *Order) float32 { return order.CORPDonation() })
}

func (o Orders) Count() int {
	return len(o)
}

func (o Orders) sum(fxn func(*Order) float32) float32 {
	var sum float32
	for _, order := range o {
		sum += fxn(order)
	}
	return sum
}

func (t *TeamOrders) All() *Orders {
	all := append(append(t.Middleton, t.MountHoreb...), t.Waunakee...)
	return &all
}
