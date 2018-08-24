package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimpelton/proc/internal/sum"

	"golang.org/x/exp/mmap"
)

func main() {
	rawFile := flag.String("f", "", "Path to raw file")

	flag.Parse()

	reader, err := mmap.Open(*rawFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Opened: ", *rawFile, "Size: ", reader.Len())

	sum := sum.SumVolume(reader)
	fmt.Println("total: ", sum)
}