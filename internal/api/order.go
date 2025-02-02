package api

import (
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type OrderType string
type Team string

const (
	TeamMoHo      = "Mount Horeb"
	TeamMiddleton = "Middleton"
	TeamWaunakee  = "Waunakee"

	TypeCamping      OrderType = "camping"
	TypeRaffle       OrderType = "raffle"
	TypeRegistration OrderType = "registration"
	TypeOther        OrderType = "other"

	riderFee float32 = 60
)

var (
	registrationRegEx = regexp.MustCompile(`\$[0-9]+`)
)

type Order struct {
	Id             string  `csv:"Order ID"`
	Subtotal       float32 `csv:"Subtotal"`
	AmountRefunded float32 `csv:"Amount Refunded"`
	Total          float32 `csv:"Total"`
	DiscountAmount float32 `csv:"Discount Amount"`
	Qty            int     `csv:"Lineitem quantity"`
	Price          float32 `csv:"Lineitem price"`
	Name           string  `csv:"Lineitem name"`
	Variant        string  `csv:"Lineitem variant"`
	Team           Team    `csv:"Product Form: Team Name"`
}

type TeamOrders struct {
	Middleton  []*Order
	MountHoreb []*Order
	Waunakee   []*Order
}

func (t *TeamOrders) All() []*Order {
	return append(append(t.Middleton, t.MountHoreb...), t.Waunakee...)
}

type Raffle struct {
	Order
	Team Team `csv:"Product Form: Team Name"`
}

func (o *Order) Type() OrderType {
	switch {
	case strings.Contains(strings.ToLower(o.Name), "registration"):
		return TypeRegistration
	case strings.Contains(strings.ToLower(o.Name), "raffle"):
		return TypeRaffle
	case strings.Contains(strings.ToLower(o.Name), "weekend"):
		return TypeCamping
	}
	return TypeOther
}

func (o *Order) IsRegistration() bool {
	return o.Type() == TypeRegistration
}

func (o *Order) RiderFeeGross() float32 {
	if !o.IsRegistration() {
		return 0
	}
	return riderFee
}
func (o *Order) RiderFeeNet() float32 {
	if !o.IsRegistration() {
		return 0
	}
	return riderFee - o.DiscountAmount - o.AmountRefunded
}

func (o *Order) Net() float32 {
	return o.Total - o.DiscountAmount - o.AmountRefunded
}

func (o *Order) BlackhawkFee() float32 {
	return o.getVariantAmount(0)
}

func (o *Order) DriftlessDonation() float32 {
	return o.getVariantAmount(1)
}

func (o *Order) CORPDonation() float32 {
	return o.getVariantAmount(2)
}

func (o *Order) getVariantAmount(index int) float32 {
	if !o.IsRegistration() {
		return 0
	}

	matches := registrationRegEx.FindAllString(o.Variant, -1)
	if len(matches) == 3 {
		amt, err := strconv.ParseFloat(matches[index][1:], 32)
		if err != nil {
			log.Errorf("Error parsing variant: %s", err)
			return 0
		}
		return float32(amt)
	}
	return 0
}
