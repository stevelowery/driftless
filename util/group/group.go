package group

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/stevelowery/driftless/internal/api"
)

func ByName(orders []*api.Order) *api.TeamOrders {
	return by(orders, func(order *api.Order) api.Team {
		switch {
		case strings.Contains(order.Name, "Waunakee"):
			return api.TeamWaunakee
		case strings.Contains(order.Name, "Middleton"):
			return api.TeamMiddleton
		case strings.Contains(order.Name, "Mount Horeb"):
			return api.TeamMoHo
		default:
			log.Fatalf("Unknown team: %s", order.Name)
			return ""
		}
	})
}

func ByTeam(orders []*api.Order) *api.TeamOrders {
	return by(orders, func(order *api.Order) api.Team {
		return order.Team
	})
}

func by(orders []*api.Order, fxn func(order *api.Order) api.Team) *api.TeamOrders {
	teamOrders := &api.TeamOrders{
		Middleton:  []*api.Order{},
		MountHoreb: []*api.Order{},
		Waunakee:   []*api.Order{},
	}
	for _, order := range orders {
		team := fxn(order)
		switch team {
		case api.TeamWaunakee:
			teamOrders.Waunakee = append(teamOrders.Waunakee, order)
		case api.TeamMiddleton:
			teamOrders.Middleton = append(teamOrders.Middleton, order)
		case api.TeamMoHo:
			teamOrders.MountHoreb = append(teamOrders.MountHoreb, order)
		}
	}
	return teamOrders
}
