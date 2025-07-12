package jobs

import (
	"fmt"
	"math"

	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/types"
)

func AllocateOrders(s storage.Storage) error {
	agents, _ := s.GetCheckedInAgents()
	orders, _ := s.GetUnassignedOrders()

	const MaxKm = 100.0
	const MaxMinutes = 600

	agentDistance := make(map[int64]float64)
	agentOrders := make(map[int64][]types.Order)

	for _, order := range orders {
		bestAgentID := int64(0)
		bestDistance := math.MaxFloat64

		for _, agent := range agents {
			if agentDistance[agent.ID] >= MaxKm {
				continue
			}

			d := Distance(order.Lat, order.Lng, 0, 0) 
			if agentDistance[agent.ID]+d > MaxKm {
				continue
			}

			if d < bestDistance {
				bestDistance = d
				bestAgentID = agent.ID
			}
		}

		if bestAgentID != 0 {
			s.AssignOrderToAgent(order.ID, bestAgentID)
			agentDistance[bestAgentID] += bestDistance
			agentOrders[bestAgentID] = append(agentOrders[bestAgentID], order)
		}
	}

	for id, list := range agentOrders {
		fmt.Printf("Agent %d assigned %d orders\n", id, len(list))
	}

	return nil
}


func Distance(lat1, lng1, lat2, lng2 float64) float64 {
	return 2.0 + math.Abs(lat1-lat2)+math.Abs(lng1-lng2)
}
