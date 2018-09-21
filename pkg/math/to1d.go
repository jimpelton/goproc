package math

func To1D(col, row, slab, maxCol, maxRow uint64) uint64 {
	return col + maxCol*(row+maxRow*slab)
}
