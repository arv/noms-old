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
		readAtLeast(r, 1<<63)
	}

	if flag.NArg() == 0 {
		size := 1 << 11
		fmt.Printf("Using random stream\n")
		readAtLeast(rand.Reader, uint64(size))
	}
}

func readAtLeast(r io.Reader, count uint64) {
	lengths := []uint64{}

	size := uint64(0)

	for i := uint64(0); i < count; i++ {
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

	a := avg(lengths)
	s := stddev(lengths)

	fmt.Printf("Size: %d, Count: %d, Avg: %d, StdDev: %d\n", size, len(lengths), a, s)
}
