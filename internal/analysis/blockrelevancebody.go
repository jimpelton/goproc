package analysis

import (
	"github.com/jimpelton/proc/pkg/indexfile"
	pmath "github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/trfunc"
	"github.com/jimpelton/proc/pkg/volume"
	log "github.com/sirupsen/logrus"
	"io"
)

type BlockRelevanceBody struct {
	Opacity  trfunc.TFOpacity
	VolStats volume.VolumeStats
	Blocks   []indexfile.FileBlock
	VDims    pmath.Vec3UI64
	BDims    pmath.Vec3UI64
	BCount   pmath.Vec3UI64

	Reader io.ReaderAt

	// NeedsNormalization is only true when the values in our data need to be
	// normalized between 0 and 1 for use in the opacity transfer function's
	// Interpolate function. If our data values are already between 0 and 1 this
	// should be false.
	NeedsNormalization bool
}

func (b *BlockRelevanceBody) F(rng Range) {
	buf := [1]byte{}

	for i := rng.Begin; i < rng.End; i = rng.Next() {
		// n, err := fillBuf(b.Reader, int64(i), buf)
		n, err := b.Reader.ReadAt(buf[:], int64(i))
		if n == 0 {
			if err != nil {
				log.Error("error reading input file:", err.Error())
			}
			break
		}

		v := float64(buf[0])
		if b.NeedsNormalization {
			v = pmath.UNorm(v, b.VolStats.Min, b.VolStats.Max)
		}
		rel := b.Opacity.Interpolate(v)

		bI := (uint64(i) % b.VDims.X()) / b.BDims.X()
		bJ := ((uint64(i) / b.VDims.X()) % b.VDims.Y()) / b.BDims.Y()
		bK := ((uint64(i) / b.VDims.X()) / b.VDims.Y()) / b.BDims.Z()

		// check block index is within block coverage
		if bI < b.BCount.X() && bJ < b.BCount.Y() && bK < b.BCount.Z() {
			bIdx := bI + b.BCount.X()*(bJ+bK*b.BCount.Y())
			b.Blocks[bIdx].Rel += rel
		}
	}
}

func (b *BlockRelevanceBody) Copy() Body {
	rval := &BlockRelevanceBody{
		Opacity:            b.Opacity,
		VolStats:           b.VolStats,
		VDims:              b.VDims,
		BDims:              b.BDims,
		BCount:             b.BCount,
		NeedsNormalization: b.NeedsNormalization,
		Reader:             b.Reader,
	}

	rval.Blocks = make([]indexfile.FileBlock, b.BCount.CompProduct())
	copy(rval.Blocks, b.Blocks)

	return rval
}

func (b *BlockRelevanceBody) Join(other ParallelReduceBody) {
	other_ := other.(*BlockRelevanceBody)
	for i, blk := range other_.Blocks {
		b.Blocks[i].Rel += blk.Rel
	}
}
