package float

import (
	"fmt"
	"strconv"
)

type Float interface {
	float32 | float64 | int | uint | uint32 | uint64
}

func TruncFloat[T float64 | float32](num T, decimal int) float64 {
	format := strconv.Itoa(decimal)
	f, _ := strconv.ParseFloat(fmt.Sprintf("%."+format+"f", num), 64)
	return f
}
