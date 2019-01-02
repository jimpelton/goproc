package analysis

import (
	"runtime"
)


func BlockLevelAnalysis(rng *Range, b *BlockRelevanceBody) *BlockRelevanceBody {
	stride := runtime.GOMAXPROCS(-1)

	bodies := make([]*BlockRelevanceBody, stride)
	wrappers := make([]*wrapper, stride)
	c := make(chan ParallelReduceBody)

	for i := range bodies {
		bodies[i] = b.Copy().(*BlockRelevanceBody)
		wrappers[i] = NewWrapper(c)
	}

	rng.stride = stride
	for tidx, body := range bodies {
		r := *rng
		r.Begin += tidx
		go wrappers[tidx].run(body, &r)
	}

	for _, w := range wrappers {
		w.waitDone()
	}

	for _, body := range bodies {
		for i, _ := range b.Blocks {
			b.Blocks[i].Rel += body.Blocks[i].Rel
		}
	}
	for _, b := range b.Blocks {
		b.Rel /= float64(b.VoxDims.CompProduct())
	}

	return b
}

type wrapper struct {
	c chan ParallelReduceBody
}

func NewWrapper(c_ chan ParallelReduceBody) *wrapper {
	return &wrapper{
		c: c_,
	}
}

func (w *wrapper) run(b ParallelReduceBody, r *Range) {
	b.F(*r)
	w.c <- b
}

func (w *wrapper) waitDone() ParallelReduceBody {
	return <-w.c
}

