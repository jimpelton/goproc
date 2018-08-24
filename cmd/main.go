package main

import (
	"flag"
	"fmt"
	"os"

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

	sumVolume(reader)

}
