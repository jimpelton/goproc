package volume

import (
	"github.com/jimpelton/mathgl/mgl64"
)

type Block struct {
	WorldDims   mgl64.Vec3 `json:"world_dims"`   // dimensions in world space
	WorldOrigin mgl64.Vec3 `json:"world_origin"` // center in world space
	Rel         float64    `json:"rel"`          // relevance value for this block
}
