package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jimpelton/mathgl/mgl64"
	"github.com/jimpelton/proc/pkg/volume"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/jimpelton/proc/internal/analysis"
	"github.com/jimpelton/proc/pkg/indexfile"
	pmath "github.com/jimpelton/proc/pkg/math"
	"github.com/jimpelton/proc/pkg/trfunc"
	"golang.org/x/exp/mmap"
)

// CmdArgs is the args for the command line.
type CmdArgs struct {
	RawFile string
	OutFile string
	TrFile  string
	DatFile string
	VolDims [3]int
	NumBlks [3]int

	Normalize bool
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

	flag.BoolVar(&args.Normalize, "norm", false, "Use normalization")
}

func main() {
	flag.Parse()
	reader, err := mmap.Open(args.RawFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.Println("Opened: ", args.RawFile, "Size: ", reader.Len())

	var opfunc *trfunc.TFOpacity
	if op, err := trfunc.OpenTFOpacityFile(args.TrFile); err != nil {
		log.Fatal("could not open tf file:", err.Error())
		os.Exit(1)
	} else {
		opfunc = op
	}

	// number of blocks along each dim of volume
	nb := pmath.Vec3UI64{uint64(args.NumBlks[0]),
		uint64(args.NumBlks[1]),
		uint64(args.NumBlks[2])}

	// volume dimensions (voxels)
	vd := pmath.Vec3UI64{uint64(args.VolDims[0]),
		uint64(args.VolDims[1]),
		uint64(args.VolDims[2])}

	// block dimensions (voxels)
	bd := vd.CompDiv(nb)

	vol := volume.Volume{
		Block: volume.Block{
			Rel: 1.0,
			WorldDims: mgl64.Vec3{1.0, 1.0, 1.0},
			WorldOrigin: mgl64.Vec3{0.0, 0.0, 0.0},
		},
		SzType: 1,
		VoxDims: vd,
	}

	log.Println("Begin volume analysis... ")
	stats := analysis.VolumeAnalysis(reader)
	log.Println("Done")

	log.Println("Begin block level analysis...")
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
			VolStats:           stats,
			Blocks:             indexfile.CreateFileBlocks(nb, vol),
			NeedsNormalization: args.Normalize,
			Reader:             reader,
		})
	log.Println("Done")

	if err := writeIndexFile(args.OutFile, args.RawFile, args.TrFile, vol, bod); err != nil {
		return
	}

	// for i, b := range bod.Blocks {
	// 	fmt.Println(i, ": ", b.Rel)
	// }
}

func writeIndexFile(outPath, rawPath, tfPath string, vol volume.Volume, body *analysis.BlockRelevanceBody) error {
	var (
		err error
		byts []byte
	)

	absRawPath, _ := filepath.Abs(rawPath)
	rawDir, rawFile := path.Split(absRawPath)
	absTfPath, _ := filepath.Abs(tfPath)
	_, tfFile := path.Split(absTfPath)

	of := indexfile.IndexFileV1{
		IndexFileHeader: indexfile.IndexFileHeader{
			TFuncName: tfFile,
			VolName: rawFile,
			VolPath: rawDir,
		},
		Volume: vol,
		Blocks: body.Blocks,
		VolStats: body.VolStats,
	}

	if byts, err = json.MarshalIndent(of, "", "  "); err != nil {
		return err
	}

	if err = ioutil.WriteFile(outPath, byts, 0644); err != nil {
		return err
	}

	return nil
}
