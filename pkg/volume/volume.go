package volume

import (
	"github.com/jimpelton/proc/pkg/math"
)

type Volume struct {
	Block
	VoxDims math.Vec3UI64 `json:"vox_dims"`
	SzType  uint64        `json:"sz_type"`
}

type VolumeStats struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Average float64 `json:"average"`
	Total   float64 `json:"total"`
}
