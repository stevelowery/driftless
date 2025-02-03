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
)

var (
	totalFxn = func(order *api.Order) float32 { return order.Net() }
)

type Options struct {
	OrdersFile     string
	RaffleFile     string
	DontationsFile string
	PayoutsFile    string
	OutputFile     string
}

type namedFunction struct {
	name string
	fxn  func(order *api.Order) float32
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
	writeLine(csv, "", api.TeamMiddleton, api.TeamMoHo, api.TeamWaunakee, "Driftless", "Total")

	orders := api.Orders{}
	raffles := api.Orders{}

	// read in the files...
	for path, slice := range map[string]any{
		opts.OrdersFile: &orders,
		opts.RaffleFile: &raffles,
	} {
		if err := io.ReadCsvInto(path, slice); err != nil {
			return err
		}
	}

	teamRegistrations := group.ByName(filter.ByType(orders, api.TypeRegistration))

	// rider fees
	riderFees := s.processRiderFees(teamRegistrations, csv)

	// raffle tickets
	teamRaffles := group.ByTeam(raffles)
	writeEmptyLine(csv)
	writeLine(csv, "Raffle", teamRaffles.Middleton.Net(), teamRaffles.MountHoreb.Net(), teamRaffles.Waunakee.Net(), "", teamRaffles.All().Net())

	// donations
	donations, err := s.processDonations(opts.DontationsFile, orders, csv)
	if err != nil {
		return err
	}

	// pass-throughs
	passThroughs := s.processPassThroughs(teamRegistrations, orders, csv)

	// total revenue
	writeEmptyLine(csv)
	writeLine(csv, "Total Revenue", "", "", "", "", riderFees+teamRaffles.All().Net()+donations+passThroughs)

	_, err = s.processPayouts(opts.PayoutsFile, csv)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) processRiderFees(teamRegistrations *api.TeamOrders, csv *csv.Writer) float32 {
	writeEmptyLine(csv)
	writeLine(csv, "Rider Count", teamRegistrations.Middleton.Count(), teamRegistrations.MountHoreb.Count(), teamRegistrations.Waunakee.Count(), "", teamRegistrations.All().Count())
	writeEmptyLine(csv)

	writeLine(csv, "Rider Fees", teamRegistrations.Middleton.RiderFeeGross(), teamRegistrations.MountHoreb.RiderFeeGross(), teamRegistrations.Waunakee.RiderFeeGross(), "", teamRegistrations.All().RiderFeeGross())
	writeLine(csv, "__Discounts", teamRegistrations.Middleton.DiscountAmount(), teamRegistrations.MountHoreb.DiscountAmount(), teamRegistrations.Waunakee.DiscountAmount(), "", teamRegistrations.All().DiscountAmount())
	writeLine(csv, "__Refunds", teamRegistrations.Middleton.AmountRefunded(), teamRegistrations.MountHoreb.AmountRefunded(), teamRegistrations.Waunakee.AmountRefunded(), "", teamRegistrations.All().AmountRefunded())
	writeLine(csv, "__Net", teamRegistrations.Middleton.RiderFeeNet(), teamRegistrations.MountHoreb.RiderFeeNet(), teamRegistrations.Waunakee.RiderFeeNet(), "", teamRegistrations.All().RiderFeeNet())

	return teamRegistrations.All().RiderFeeNet()
}

func (s *service) processDonations(path string, orders api.Orders, csv *csv.Writer) (float32, error) {

	teamDonations := group.ByName(filter.ByType(orders, api.TypeRegistration))

	donations := []*api.Donation{}
	if err := io.ReadCsvInto(path, &donations); err != nil {
		return 0, err
	}

	var driftlessTotal float32
	for _, donation := range donations {
		driftlessTotal += donation.Total
	}

	teamsTotal := teamDonations.All().DriftlessDonation()

	writeEmptyLine(csv)
	writeLine(csv, "Donations", "", "", "", "", "")
	writeLine(csv, "__Standalone", "", "", "", fmt.Sprint(driftlessTotal), fmt.Sprint(driftlessTotal))
	writeLine(csv, "__With Registration", teamDonations.Middleton.DriftlessDonation(), teamDonations.MountHoreb.DriftlessDonation(), teamDonations.Waunakee.DriftlessDonation(), "", teamsTotal)
	writeLine(csv, "__Subtotal", "", "", "", "", driftlessTotal+teamsTotal)

	return driftlessTotal + teamsTotal, nil
}

func (s *service) processPassThroughs(teamorders *api.TeamOrders, orders api.Orders, csv *csv.Writer) float32 {
	camping := filter.ByType(orders, api.TypeCamping).Net()
	subtotal := teamorders.All().BlackhawkFee() + teamorders.All().CORPDonation() + camping

	writeEmptyLine(csv)
	writeLine(csv, "Pass-throughs", "", "", "", "", "")
	writeLine(csv, "__Blackhawk Fee", teamorders.Middleton.BlackhawkFee(), teamorders.MountHoreb.BlackhawkFee(), teamorders.Waunakee.BlackhawkFee(), "", teamorders.All().BlackhawkFee())
	writeLine(csv, "__CORP Donation", teamorders.Middleton.CORPDonation(), teamorders.MountHoreb.CORPDonation(), teamorders.Waunakee.CORPDonation(), "", teamorders.All().CORPDonation())
	writeLine(csv, "__Camping", "", "", "", camping, camping)
	writeLine(csv, "__Subtotal", "", "", "", "", subtotal)

	return subtotal
}

func (s *service) processPayouts(path string, csv *csv.Writer) (float32, error) {

	writeEmptyLine(csv)
	writeLine(csv, "Payouts", "", "", "", "", "")

	payouts := api.Payouts{}
	if err := io.ReadCsvInto(path, &payouts); err != nil {
		return 0, err
	}

	totals := payouts.Total()
	writeLine(csv, "__Charges", "", "", "", totals.Charges, totals.Charges)
	writeLine(csv, "__Refunds", "", "", "", totals.Refunds, totals.Refunds)
	writeLine(csv, "__Fees", "", "", "", totals.Fees, totals.Fees)
	writeLine(csv, "__Net", "", "", "", totals.Net, totals.Net)

	return totals.Net, nil
}

func writeEmptyLine(csv *csv.Writer) error {
	return csv.Write([]string{"", "", "", "", "", ""})
}

func writeLine(csv *csv.Writer, values ...any) error {
	vals := []string{}
	for _, value := range values {
		vals = append(vals, fmt.Sprint(value))
	}
	return csv.Write(vals)
}
