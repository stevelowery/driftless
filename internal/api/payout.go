package api

type Payouts []*Payout

type Payout struct {
	Charges float32 `csv:"Charges"`
	Refunds float32 `csv:"Refunds"`
	Fees    float32 `csv:"Fees"`
	Net     float32 `csv:"Net Amount"`
}

func (p *Payouts) Total() *Payout {
	var total Payout
	for _, payout := range *p {
		total.Charges += payout.Charges
		total.Refunds += payout.Refunds
		total.Fees += payout.Fees
		total.Net += payout.Net
	}
	return &total
}
