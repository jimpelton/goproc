package analysis

import (
	// log "github.com/sirupsen/logrus"
	"runtime"
)


func BlockAnalysis(rng *Range, b Body) Body {
	stride := runtime.GOMAXPROCS(-1)
	bodies := make([]Body, stride)
	for i := range bodies {
		bodies[i] = b.Copy()
	}

	rng.stride = stride
	for tidx, b := range bodies {
		r := *rng
		r.Begin += tidx
		go b.F(*rng)
	}

	// combine results into b
	// for _, b := range bodies {
	//
	// }

	return b
}

