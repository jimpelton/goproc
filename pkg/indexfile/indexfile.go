package indexfile

import (
	"github.com/jimpelton/mathgl/mgl64"
	"github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/volume"
)

type IndexFileHeader struct {
	Magic       uint16    `json:"magic"`
	FileVersion uint16    `json:"file_version"`
	VolName     [256]byte `json:"vol_name"`
	VolPath     [512]byte `json:"vol_path"`
	TFuncName   [256]byte `json:"t_func_name"`
}

type IndexFileV1 struct {
	IndexFileHeader
	Volume volume.Volume  `json:"volume"`
	Blocks []volume.Block `json:"blocks"`
}

func NewIndexFileV1() *IndexFileV1 {
	return &IndexFileV1{
		IndexFileHeader: IndexFileHeader{
			Magic:       7376,
			FileVersion: 1,
			VolName:     [256]byte{},
			VolPath:     [512]byte{},
			TFuncName:   [256]byte{},
		},
		Volume: volume.Volume{},
		Blocks: []volume.Block{},
	}
}

type FileBlock struct {
	volume.Block
	VoxDims math.Vec3UI64 `json:"dims"`   // voxel dimensions
	Ijk     math.Vec3UI64 `json:"ijk"`    // block 3D index
	Offset  uint64        `json:"offset"` // byte offset into file
}

// CreateVolumeBlocks creates blocks in world space within the given Volume
func CreateFileBlocks(nblocks math.Vec3UI64, vol volume.Volume) (blocks []*FileBlock) {
	blkDimsWld := vol.WorldDims.CompDiv(
		mgl64.Vec3{float64(nblocks.X()), float64(nblocks.Y()), float64(nblocks.Z())})

	blkDimsVox := vol.VoxDims.CompDiv(nblocks)

	for k := uint64(0); k < nblocks.Z(); k++ {
		for j := uint64(0); j < nblocks.Y(); j++ {
			for i := uint64(0); i < nblocks.X(); i++ {
				worldLoc := mgl64.Vec3{
					vol.WorldDims.X()*float64(i) - 0.5,
					vol.WorldDims.Y()*float64(j) - 0.5,
					vol.WorldDims.Z()*float64(k) - 0.5}
				// origin = worldLoc + (worldLoc + wd)) * 0.5 (midpt formula)
				origin := worldLoc.Add(blkDimsWld).Add(worldLoc).Mul(0.5)
				startVoxel := math.Vec3UI64{i, j, k}.CompMul(blkDimsVox)

				blocks = append(blocks,
					&FileBlock{
						Block: volume.Block{
							WorldOrigin: origin,
							WorldDims:   blkDimsWld,
							Rel:         0.0,
						},
						Ijk:     math.Vec3UI64{i, j, k},
						VoxDims: blkDimsVox,
						Offset: math.To1D(startVoxel.X(), startVoxel.Y(), startVoxel.Z(),
							vol.VoxDims.X(), vol.VoxDims.Y()) * vol.SzType,
					})
			}
		}
	}
	return nil
}
