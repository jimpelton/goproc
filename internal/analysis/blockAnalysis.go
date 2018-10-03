package analysis

import (
	"github.com/jimpelton/proc/internal/result"
	"github.com/jimpelton/proc/pkg/indexfile"
	pmath "github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/trfunc"
	"github.com/jimpelton/proc/pkg/volume"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/mmap"
	"math"
	"runtime"
)

type AnalyzedBlocks struct {
	Blocks []*indexfile.FileBlock
}

//type Func func(buffer []byte, tidx, stride int) interface{}

type Body interface {
	F(buffer []byte, tidx, stride, end int) interface{}
}

type Args struct {
	NBlocks pmath.Vec3UI64
	Volume  volume.Volume
	Body    Body
}

func BlockAnalysis(r *mmap.ReaderAt, args Args) AnalyzedBlocks {
	np := runtime.GOMAXPROCS(-1)
	results := make([]*result.Result, np)
	for i, _ := range results {
		results[i] = result.NewResult()
		start := 0
		end := r.Len()
		go runAnalysis(r, i, start, end, np, results[i], args)
	}

	// for _, res := range results {

	// }

	return AnalyzedBlocks{}
}

func runAnalysis(r *mmap.ReaderAt, tidx, start, end, stride int, res *result.Result, args Args) {
	bs := AnalyzedBlocks{Blocks: indexfile.CreateFileBlocks(args.NBlocks, args.Volume)}
	var (
		buf []byte
	)
	buf = make([]byte, int(math.Pow(2, 15)))
	for start < end {
		n, err := fillBuf(r, int64(start), buf)
		if n == 0 {
			if err != nil {
				log.Error("error reading input file:", err.Error())
			}
			break
		}
		start += n
		val := args.Body.F(buf, tidx, stride, end)
	}
}

type BlockRelevanceBody struct {
	opacity  trfunc.TFOpacity
	volStats volume.VolumeStats
	// needsNormalization is only true when the values in our data need to be
	// normalized between 0 and 1 for use in the opacity transfer function's
	// Interpolate function. If our data values are already between 0 and 1 this
	// should be false.
	needsNormalization bool
}

func (b *BlockRelevanceBody) F(buf []byte, tidx, stride, end int) interface{} {
	for i := tidx; i < end; i += stride {
		v := float64(buf[i])
		if b.needsNormalization {
			v = UNorm(v, b.volStats.Min, b.volStats.Max)
		}
	}
	return nil
}

// UNorm performs unity-based normalization.
//
// Scales v to within [0 .. 1]
func UNorm(v, min, max float64) float64 {
	return (v - min) / (max - min)
}
