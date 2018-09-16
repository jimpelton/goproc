package volume

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Block struct {
	Dims        [3]uint64  `json:"dims"`         // voxel dimensions
	Ijk         [3]uint64  `json:"ijk"`          // block 3D index
	WorldDims   mgl64.Vec3 `json:"world_dims"`   // dimensions in world space
	WorldOrigin mgl64.Vec3 `json:"world_origin"` // center in world space
	Rel         float64    `json:"rel"`          // relevance value for this block
}
