package analysis

import (
	// log "github.com/sirupsen/logrus"
	"math"
	"runtime"
)


func ParallelReduce(rng *Range, b ParallelReduceBody) Body {
	stride := runtime.GOMAXPROCS(-1)

	bodies := make([]ParallelReduceBody, stride)
	wrappers := make([]*wrapper, stride)
	c := make(chan ParallelReduceBody)
	for i := range bodies {
		bodies[i] = b.Copy().(ParallelReduceBody)
		wrappers[i] = NewWrapper(c)
	}

	rng.stride = stride
	for tidx, b := range bodies {
		// copy the range
		r := *rng
		r.Begin += tidx
		go wrappers[tidx].run(b, &r)
	}

	for i := 0; i <= len(bodies); {
		b1 := <-c
		i++
		b2 := <-c
		i++
		b1.Join(b2)
	}

	levels := math.Log2(float64(len(bodies)))
	for lev := levels; lev >= 0; lev-- {
		// combine results into b
		for i := 0; i < len(bodies); i += 2 {

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

// func (w *wrapper) waitDone() {
// 	w.wait.Wait()
// }

