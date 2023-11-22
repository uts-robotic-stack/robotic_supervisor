package device

import (
	"math/rand"

	"github.com/containrrr/watchtower/pkg/types"
)

func GetHardwareStatus() (types.HardwareStatus, error) {
	return types.HardwareStatus{
		Cpu:         rand.Float64() * 100.0,
		Temperature: rand.Float64() * 100.0,
		Ram:         rand.Float64() * 100.0,
		Storage:     rand.Float64() * 100.0,
	}, nil
}
