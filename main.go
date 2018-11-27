package main

import (
	"flag"
	"fmt"
	"github.com/jimpelton/proc/internal/analysis"
	"github.com/jimpelton/proc/pkg/indexfile"
	pmath "github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/trfunc"
	"github.com/jimpelton/proc/pkg/volume"
	"golang.org/x/exp/mmap"
	"log"
	"os"
)

type CmdArgs struct {
	RawFile string
	OutFile string
	TrFile  string
	DatFile string
	VolDims [3]int
	NumBlks [3]int
}

var (
	args CmdArgs
)

func init() {
	flag.StringVar(&args.RawFile, "f", "", "Path to raw file")
	flag.StringVar(&args.OutFile, "o", "", "Output index file")
	flag.StringVar(&args.TrFile, "t", "", "Transfer function file")
	flag.StringVar(&args.DatFile, "d", "", "Dat file (currently ignored)")

	flag.IntVar(&args.VolDims[0], "vx", 1, "Volume X-dim")
	flag.IntVar(&args.VolDims[1], "vy", 1, "Volume Y-dim")
	flag.IntVar(&args.VolDims[2], "vz", 1, "Volume Z-dim")

	flag.IntVar(&args.NumBlks[0], "bx", 1, "Num blocks X-dim")
	flag.IntVar(&args.NumBlks[1], "by", 1, "Num blocks Y-dim")
	flag.IntVar(&args.NumBlks[2], "bz", 1, "Num blocks Z-dim")
}

func main() {
	flag.Parse()
	reader, err := mmap.Open(args.RawFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Opened: ", args.RawFile, "Size: ", reader.Len())
	sumtotal := analysis.VolumeAnalysis(reader)
	fmt.Printf("total: %f\n", sumtotal)

	var opfunc *trfunc.TFOpacity
	if op, err := trfunc.OpenTFOpacityFile(args.TrFile); err != nil {
		log.Fatal("could not open tf file:", err.Error())
		os.Exit(1)
	} else {
		opfunc = op
	}

	nb := pmath.Vec3UI64{uint64(args.NumBlks[0]),
		uint64(args.NumBlks[1]),
		uint64(args.NumBlks[2])}
	vd := pmath.Vec3UI64{uint64(args.VolDims[0]),
		uint64(args.VolDims[1]),
		uint64(args.VolDims[2])}
	bd := vd.CompDiv(nb)

	bod := analysis.BlockLevelAnalysis(
		&analysis.Range{
			Begin: 0,
			End:   reader.Len(),
		},
		&analysis.BlockRelevanceBody{
			Opacity:            *opfunc,
			BCount:             nb,
			BDims:              bd,
			VDims:              vd,
			VolStats:           volume.VolumeStats{},
			Blocks:             make([]indexfile.FileBlock, nb.CompProduct()),
			NeedsNormalization: false,
			Reader:             reader,
		})


}
