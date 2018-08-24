package main

import (
	"runtime"

	"golang.org/x/exp/mmap"
)

func sumVolume(r *mmap.ReaderAt) {
	np := runtime.GOMAXPROCS(-1)
	splits := r.Len() / np

	results := make([]*Result, splits)
	for i, _ := range results {
		results[i] = NewResult()
		start := i * splits
		end := start + splits
		go sum(r, start, end, results[i])
	}

	sum := float64(0.0)
	for _, res := range results {
		sum += res.Wait().(float64)
	}

}

func sum(r *mmap.ReaderAt, start int, end int, res *Result) {
	var (
		buf []byte
		sum float64
		//err error
	)
	buf = make([]byte, 1024)
	for start < end {
		n := fillBuf(r, int64(start), buf)
		start += n
		for i := 0; i < n; i++ {
			sum += float64(buf[i])
		}
	}
	res.Done(sum)
}

func fillBuf(r *mmap.ReaderAt, off int64, buf []byte) int {
	n := 0
	bufLen := len(buf)
	for n < bufLen {
		i, err := r.ReadAt(buf[n:], off)
		n += i
		if err != nil {
			break
		}
	}
	return n
}
