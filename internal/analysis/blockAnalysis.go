package analysis

import (
	// log "github.com/sirupsen/logrus"
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
	for tidx, b := range bodies {
		r := *rng
		r.Begin += tidx
		go wrappers[tidx].run(b, &r)
	}

	for _, w := range wrappers {
		w.waitDone()
	}

	for _, bdone := range bodies {
		for i, _ := range b.Blocks {
			b.Blocks[i].Rel += bdone.Blocks[i].Rel
		}
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

