package volume

type Volume struct {
	Block
}

type VolumeStats struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Average float64 `json:"average"`
	Total   float64 `json:"total"`
}


func CreateVolumeBlocks(nblocks [3]uint64) []Block {
	return nil
}