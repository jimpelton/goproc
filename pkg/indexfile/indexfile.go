package indexfile

import (
	"github.com/jimpelton/proc/pkg/volume"
)

type IndexFileHeader struct {
	Magic       uint16    `json:"magic"`
	FileVersion uint16    `json:"file_version"`
	Name        [256]byte `json:"name"`
	Path        [512]byte `json:"path"`
	VolName     [256]byte `json:"vol_name"`
	VolPath     [512]byte `json:"vol_path"`
}

type IndexFileV1 struct {
	IndexFileHeader
	Volume volume.Volume  `json:"volume"`
	Blocks []volume.Block `json:"blocks"`
}

type FileBlock struct {
	Block      volume.Block `json:"block"`
	FileOffset uint64       `json:"file_offset"`
}