package analysis

import (
	"fmt"
	"math"
	"runtime"

	"github.com/jimpelton/proc/internal/result"
	"github.com/jimpelton/proc/pkg/volume"
	"golang.org/x/exp/mmap"

	log "github.com/sirupsen/logrus"
)

// VolumeAnalysis computes the sum of elements in the volume and min, max, avg values.
func VolumeAnalysis(r *mmap.ReaderAt) (stats volume.VolumeStats) {
	np := runtime.GOMAXPROCS(-1)
	results := make([]*result.Result, np)
	for i, _ := range results {
		results[i] = result.NewResult()
		start := 0
		end := r.Len()
		go runSum(r, i, start, end, np, results[i] /*, progress*/)
	}

	stats.Min = math.MaxFloat64
	stats.Max = -math.MaxFloat64

	for _, res := range results {
		s := res.Wait().(volume.VolumeStats)

		stats.Total += s.Total

		if s.Min < stats.Min {
			stats.Min = s.Min
		}
		if s.Max > stats.Max {
			stats.Max = s.Max
		}
	}

	stats.Average = stats.Total / float64(r.Len())

	return
}

func runSum(r *mmap.ReaderAt,
	tidx, start, end, stride int, res *result.Result /*progress chan int*/) {
	var (
		buf []byte
		sum float64
		min float64 = math.MaxFloat64
		max float64 = -math.MaxFloat64
	)
	buf = make([]byte, int(math.Pow(2, 15)))
	for start < end {
		n, err := fillBuf(r, int64(start), buf)
		if n == 0 {
			if err != nil {
				log.Error("error:", err.Error())
			}
			break
		}
		start += n

		for i := tidx; i < n; i += stride {
			v := float64(buf[i])
			sum += v

			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}

		if err != nil {
			fmt.Println("error:", err.Error())
			break
		}
	}
	res.Done(volume.VolumeStats{
		Min:     min,
		Max:     max,
		Total:   sum,
		Average: 0,
	})
	fmt.Println("Done.")
}

func fillBuf(r *mmap.ReaderAt, off int64, buf []byte) (int, error) {
	n := 0
	bufLen := len(buf)
	for n < bufLen {
		i, err := r.ReadAt(buf[n:], off)
		n += i
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
