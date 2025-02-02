package squarespace

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/stevelowery/driftless/internal/api"
	"github.com/stevelowery/driftless/util/filter"
	"github.com/stevelowery/driftless/util/group"
	"github.com/stevelowery/driftless/util/io"
	"github.com/stevelowery/driftless/util/math"
)

var (
	totalFxn = func(order *api.Order) float32 { return order.Net() }
)

type Options struct {
	OrdersFile     string
	RaffleFile     string
	DontationsFile string
	OutputFile     string
}

type Service interface {
	Process(ctx context.Context, opts *Options) error
}

type service struct{}

func New() Service {
	return &service{}
}

func (s *service) Process(ctx context.Context, opts *Options) error {
	log.Info("Processing orders...")
	defer log.Info("Processing complete.")

	// create the report...
	output, err := os.Create(opts.OutputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	csv := csv.NewWriter(output)
	defer csv.Flush()

	// header
	csv.Write([]string{"", api.TeamMiddleton, api.TeamMoHo, api.TeamWaunakee, "Driftless", "Total"})

	orders := []*api.Order{}
	raffles := []*api.Order{}
	donations := []*api.Donation{}

	// read in the files...
	for path, slice := range map[string]any{
		opts.OrdersFile:     &orders,
		opts.RaffleFile:     &raffles,
		opts.DontationsFile: &donations,
	} {
		if err := io.ReadCsvInto(path, slice); err != nil {
			return err
		}
	}

	teamRegistrations := group.ByName(filter.ByType(orders, api.TypeRegistration))
	teamRaffles := group.ByTeam(raffles)

	csv.Write([]string{"", "", "", "", "", ""})
	csv.Write([]string{"Rider Count", math.Count(teamRegistrations.Middleton), math.Count(teamRegistrations.MountHoreb), math.Count(teamRegistrations.Waunakee), "", math.Count(teamRegistrations.All())})
	csv.Write([]string{"", "", "", "", "", ""})

	csv.Write([]string{"Revenue", "", "", "", "", ""})
	for name, fxn := range map[string]func(order *api.Order) float32{
		"__Rider Fees":  func(order *api.Order) float32 { return order.RiderFeeGross() },
		"____Discounts": func(order *api.Order) float32 { return order.DiscountAmount },
		"____Refunds":   func(order *api.Order) float32 { return order.AmountRefunded },
		"____Net":       func(order *api.Order) float32 { return order.RiderFeeNet() },
	} {
		if err := s.reportRegistrationData(teamRegistrations, csv, name, fxn); err != nil {
			return err
		}
	}

	csv.Write([]string{"", "", "", "", "", ""})
	csv.Write([]string{"__Raffle", math.Sum(teamRaffles.Middleton, totalFxn), math.Sum(teamRaffles.MountHoreb, totalFxn), math.Sum(teamRaffles.Waunakee, totalFxn), "", math.Sum(teamRaffles.All(), totalFxn)})

	if err := s.processDonations(orders, donations, csv); err != nil {
		return err
	}

	csv.Write([]string{"", "", "", "", "", ""})
	csv.Write([]string{"Pass-throughs", "", "", "", "", ""})
	for name, fxn := range map[string]func(order *api.Order) float32{
		"__Blackhawk Fee": func(order *api.Order) float32 { return order.BlackhawkFee() },
		"__CORP Donation": func(order *api.Order) float32 { return order.CORPDonation() },
	} {
		if err := s.reportRegistrationData(teamRegistrations, csv, name, fxn); err != nil {
			return err
		}
	}

	campingTotal := math.Sum(filter.ByType(orders, api.TypeCamping), totalFxn)
	csv.Write([]string{"__Camping", "", "", "", campingTotal, campingTotal})

	return nil
}

func (s *service) reportRiderFees(teamOrders *api.TeamOrders, csv *csv.Writer) error {

	for name, fxn := range map[string]func(order *api.Order) float32{
		"__Rider Fees": func(order *api.Order) float32 { return order.RiderFeeGross() },
		"__Discounts":  func(order *api.Order) float32 { return order.DiscountAmount },
		"__Refunds":    func(order *api.Order) float32 { return order.AmountRefunded },
		"__Net":        func(order *api.Order) float32 { return order.RiderFeeNet() },
	} {
		if err := s.reportRegistrationData(teamOrders, csv, name, fxn); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) reportRegistrationData(teamOrders *api.TeamOrders, csv *csv.Writer, name string, fxn func(order *api.Order) float32) error {
	if err := csv.Write([]string{name, math.Sum(teamOrders.Middleton, fxn), math.Sum(teamOrders.MountHoreb, fxn), math.Sum(teamOrders.Waunakee, fxn), "", math.Sum(teamOrders.All(), fxn)}); err != nil {
		return err
	}
	return nil
}

func (s *service) processDonations(orders []*api.Order, donations []*api.Donation, csv *csv.Writer) error {
	teamDonations := group.ByName(filter.ByType(orders, api.TypeRegistration))

	var driftlessTotal float32
	for _, donation := range donations {
		driftlessTotal += donation.Total
	}

	csv.Write([]string{"", "", "", "", "", ""})
	donationFxn := func(order *api.Order) float32 { return order.DriftlessDonation() }
	csv.Write([]string{"__Donations", "", "", "", "", ""})
	csv.Write([]string{"____Standalone", "", "", "", fmt.Sprint(driftlessTotal), fmt.Sprint(driftlessTotal)})
	csv.Write([]string{"____With Registration", math.Sum(teamDonations.Middleton, donationFxn), math.Sum(teamDonations.MountHoreb, donationFxn), math.Sum(teamDonations.Waunakee, donationFxn), "", math.Sum(teamDonations.All(), donationFxn)})

	return nil
}

func (s *service) processCamping(orders []*api.Order, csv *csv.Writer) error {
	camping := filter.ByType(orders, api.TypeCamping)
	csv.Write([]string{"__Camping", "", "", "", math.Sum(camping, totalFxn), ""})
	return nil
}
