package math

import (
	"fmt"
	"math"

	"github.com/stevelowery/driftless/internal/api"
)

func Count(orders []*api.Order) string {
	return fmt.Sprint(len(orders))
}

func Sum(orders []*api.Order, fxn func(*api.Order) float32) string {
	var sum float32
	for _, order := range orders {
		sum += fxn(order)
	}

	return fmt.Sprint(float32(math.Round(float64(sum)*100) / 100))
}
