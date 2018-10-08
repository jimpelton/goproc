package analysis

import (
	"github.com/jimpelton/proc/pkg/indexfile"
	pmath "github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/trfunc"
	"github.com/jimpelton/proc/pkg/volume"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/mmap"
	"golang.org/x/text/encoding/internal/identifier"
	"io"
	"runtime"
)

type Body interface {
	// F is called by each goroutine and passed in the shared buffer.
	F(Range)
	// A Body must be able to copy itself -- each goroutine is given a copy
	// of the body.
	Copy() Body
}

// type Range interface {
// 	At() int
// 	Next() int
// 	Begin() int
// 	End() int
//
// 	setStride(int)
// }

type Range struct {
	End int
	Begin int
	stride int
	next int
}


func (l *Range) Next() (n int) {
	n = l.next
	l.next += l.stride
	return
}

func BlockAnalysis(rng *Range, b Body) Body {
	stride := runtime.GOMAXPROCS(-1)
	bodies := make([]Body, stride)
	for i := range bodies {
		bodies[i] = b.Copy()
	}

	rng.stride = stride
	for i, b := range bodies {
		r := *rng
		r.Begin += tidx
		go b.F(rng)
	}

	// combine results into b
	// for _, b := range bodies {
	//
	// }

	return b
}

// runAnalysis reads r and fills a buffer with data. That buffer is processed by the Body.
// func runAnalysis(tidx int, rng Range, b Body) {
//
// 	for rng.Begin < rng.End {
//
// 		n, err := fillBuf(r, int64(start), buf)
//
// 		if n == 0 {
// 			if err != nil {
// 				log.Error("error reading input file:", err.Error())
// 			}
// 			break
// 		}
//
// 	}
//
// }

type BlockRelevanceBody struct {
	Opacity  trfunc.TFOpacity
	VolStats volume.VolumeStats
	Blocks   []indexfile.FileBlock
	VDims    pmath.Vec3UI64
	BDims    pmath.Vec3UI64
	BCount   pmath.Vec3UI64

	reader io.ReaderAt

	// needsNormalization is only true when the values in our data need to be
	// normalized between 0 and 1 for use in the opacity transfer function's
	// Interpolate function. If our data values are already between 0 and 1 this
	// should be false.
	needsNormalization bool
}

func (b *BlockRelevanceBody) F(rng Range) {
	buf := [1]byte{}

	for i := rng.Begin; i < rng.End; i = rng.Next() {
		// n, err := fillBuf(b.reader, int64(i), buf)
		n, err := b.reader.ReadAt(buf[:], int64(i))
		if n == 0 {
			if err != nil {
				log.Error("error reading input file:", err.Error())
			}
			break
		}

		v := float64(buf[0])
		if b.needsNormalization {
			v = UNorm(v, b.VolStats.Min, b.VolStats.Max)
		}
		rel := b.Opacity.Interpolate(v)

		bI := (uint64(i) % b.VDims.X()) / b.BDims.X()
		bJ := ((uint64(i) / b.VDims.X()) % b.VDims.Y()) / b.BDims.Y()
		bK := ((uint64(i) / b.VDims.X()) / b.VDims.Y()) / b.BDims.Z()

		// check block index is within block coverage
		if bI < b.BCount.X() && bJ < b.BCount.Y() && bK < b.BCount.Z() {
			bIdx := bI + b.BCount.X() * (bJ + bK * b.BCount.Y())
			b.Blocks[bIdx].Rel += rel
		}
	}
}


// UNorm performs unity-based normalization.
//
// Scales v to within [0 .. 1]
func UNorm(v, min, max float64) float64 {
	return (v - min) / (max - min)
}
