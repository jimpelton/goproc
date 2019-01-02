package indexfile

import (
	"github.com/jimpelton/mathgl/mgl64"
	"github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/volume"
)

const magic = 7376
const fileversion = 1

type IndexFileHeader struct {
	VolName   string `json:"vol_name"`
	VolPath   string `json:"vol_path"`
	TFuncName string `json:"t_func_name"`
}

type IndexFileV1 struct {
	IndexFileHeader
	Volume   volume.Volume      `json:"volume"`
	VolStats volume.VolumeStats `json:"vol_stats"`
	Blocks   []FileBlock        `json:"blocks"`
}

// func NewIndexFileV1() *IndexFileV1 {
// 	return &IndexFileV1{
// 		IndexFileHeader: IndexFileHeader{
// 		},
// 		Volume: volume.Volume{},
// 		Blocks: []volume.Block{},
// 	}
// }

// A Block entry in the IndexFile.
type FileBlock struct {
	volume.Block
	Index   uint64        `json:"index"`
	VoxDims math.Vec3UI64 `json:"dims"`   // voxel dimensions
	Ijk     math.Vec3UI64 `json:"ijk"`    // block 3D index
	Offset  uint64        `json:"offset"` // byte offset into file
}

// CreateFileBlocks creates blocks in world space within the given Volume
func CreateFileBlocks(nBlocks math.Vec3UI64, vol volume.Volume) (blocks []FileBlock) {
	// the size of each Block in world dimensions.
	blkDimsWld := vol.WorldDims.CompDiv(
		mgl64.Vec3{float64(nBlocks.X()), float64(nBlocks.Y()), float64(nBlocks.Z())})

	// the voxel length of each dimension
	blkDimsVox := vol.VoxDims.CompDiv(nBlocks)

	blocks = make([]FileBlock, nBlocks.CompProduct())
	index := uint64(0)

	for k := uint64(0); k < nBlocks.Z(); k++ {
		for j := uint64(0); j < nBlocks.Y(); j++ {
			for i := uint64(0); i < nBlocks.X(); i++ {
				worldLoc := mgl64.Vec3{
					vol.WorldDims.X()*float64(i) - (vol.WorldDims.X() * 0.5),
					vol.WorldDims.Y()*float64(j) - (vol.WorldDims.Y() * 0.5),
					vol.WorldDims.Z()*float64(k) - (vol.WorldDims.Z() * 0.5)}

				// origin = worldLoc + (worldLoc + wd)) * 0.5 (midpt formula)
				origin := worldLoc.Add(blkDimsWld).Add(worldLoc).Mul(0.5)
				startVoxel := math.Vec3UI64{i, j, k}.CompMul(blkDimsVox)

				blocks[index] = FileBlock{
					Block: volume.Block{
						WorldOrigin: origin,
						WorldDims:   blkDimsWld,
						Rel:         0.0,
					},
					Index:   index,
					Ijk:     math.Vec3UI64{i, j, k},
					VoxDims: blkDimsVox,
					Offset: math.To1D(startVoxel.X(), startVoxel.Y(), startVoxel.Z(),
						vol.VoxDims.X(), vol.VoxDims.Y()) * vol.SzType,
				}
				index++
			}
		}
	}
	return blocks
}
