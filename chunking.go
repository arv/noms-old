package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/attic-labs/noms/types"
)

func avg(ns []uint64) uint64 {
	sum := uint64(0)
	for _, n := range ns {
		sum += n
	}
	return sum / uint64(len(ns))
}

func stddev(ns []uint64) uint64 {
	a := avg(ns)

	deviations := make([]uint64, len(ns))
	for i, n := range ns {
		deviations[i] = uint64(math.Pow(float64(n-a), 2))
	}
	variance := avg(deviations)
	return uint64(math.Sqrt(float64(variance)))
}

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		fmt.Printf("File name: %s\n", arg)
		f, err := os.Open(arg)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		r := bufio.NewReader(f)
		readData(r)
	}

	if flag.NArg() == 0 {
		fmt.Printf("Using random stream\n")
		readData(io.LimitReader(rand.Reader, 10e6))
	}
}

func readData(r io.Reader) {
	lengths := []uint64{}

	size := uint64(0)

	for {
		w := bytes.Buffer{}
		n, err := types.CopyChunk(&w, r)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		size += n
		lengths = append(lengths, n)
	}

	var a, s uint64
	// skip the last one since it will skew the average/stddev.
	if len(lengths) > 2 {
		a = avg(lengths[:len(lengths)-1])
		s = stddev(lengths[:len(lengths)-1])
	}

	fmt.Printf("Size: %d, Count: %d, Avg: %d, StdDev: %d\n", size, len(lengths), a, s)
}
