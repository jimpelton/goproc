package volume

import (
	"github.com/jimpelton/mathgl/mgl64"
	"github.com/jimpelton/proc/pkg/math"
)

type Volume struct {
	Block
}

type VolumeStats struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Average float64 `json:"average"`
	Total   float64 `json:"total"`
}


func to1D(col, row, slab, maxCol, maxRow uint64) uint64 {
	return col + maxCol * (row + maxRow * slab)
}

// CreateVolumeBlocks creates blocks in world space within the given Volume
func CreateVolumeBlocks(nblocks math.Vec3UI64, vol Volume) (blocks []Block) {
	blocks = make([]Block, nblocks.CProduct())

	wldDims := vol.WorldDims.CompDiv(mgl64.Vec3{float64(nblocks.X()), float64(nblocks.Y()), float64(nblocks.Z())})

	//		mgl64.Vec3{
	//	vol.WorldDims.X() / float64(nblocks.X()),
	//	vol.WorldDims.Y() / float64(nblocks.Y()),
	//	vol.WorldDims.Z() / float64(nblocks.Z()),
	//}

	// block voxel dimensions
	// volume vox dimensions / num blocks
	//voxDims := math.Vec3UI64{
	//	vol.VoxDims.X() / nblocks.X()
	//}

	for k := uint64(0); k < nblocks.Z(); k++ {
		for j := uint64(0); j < nblocks.Y(); j++ {
			for i := uint64(0); i < nblocks.X(); i++ {
				idx := to1D(i, j, k, nblocks.X(), nblocks.Y())
				worldLoc := mgl64.Vec3{
					vol.WorldDims.X() * float64(i) - 0.5,
					vol.WorldDims.Y() * float64(j) - 0.5,
					vol.WorldDims.Z() * float64(k) - 0.5 }

				// midpoint formula
				// origin = worldLoc + (worldLoc + wd)) * 0.5
				origin := worldLoc.Add(wldDims).Add(worldLoc).Mul(0.5)



			}
		}
	}
	return nil
}