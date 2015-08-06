package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"github.com/attic-labs/noms/types"

	"github.com/attic-labs/buzhash"
	"github.com/attic-labs/noms/adler32"
	"github.com/attic-labs/noms/rabinkarp"
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

func BuzHash(windowSize uint32) hash.Hash32 {
	return buzhash.NewBuzHash(windowSize)
}

func RabinKarp(windowSize uint32) hash.Hash32 {
	return rabinkarp.NewRabinKarp(windowSize)
}

func Adler32(windowSize uint32) hash.Hash32 {
	return adler32.NewAdler32(windowSize)
}

type test struct {
	name    string
	pattern string
}

type suite struct {
	name string
	f    func(ws uint32) hash.Hash32
}

func main() {
	flag.Parse()

	tests := []test{
		{"Movie", "*.mkv"},
		{"Jpegs", "*.jpg"},
		{"Pdfs", "*.pdf"},
	}

	suites := []suite{
		{"Adler32", Adler32},
		{"BuzHash", BuzHash},
		{"RabinKarp", RabinKarp},
	}

	buf, _ := ioutil.ReadAll(io.LimitReader(rand.Reader, 10e6))

	fmt.Printf("Suite, Test, Size, Count, Average, StdDev\n")
	for _, s := range suites {
		types.HashFunc = s.f
		readData(s.name, "Stream", bytes.NewBuffer(buf))
		for _, t := range tests {
			doFileTest(s.name, t.name, flag.Arg(0)+t.pattern, s.f)
		}
	}
}

func readData(suite, test string, r io.Reader) {
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

	fmt.Printf("%s, %s, %d, %d, %d, %d\n", suite, test, size, len(lengths), a, s)
}

func doFileTest(suite, test, pattern string, f func(ws uint32) hash.Hash32) {
	matches, _ := filepath.Glob(pattern)

	for _, n := range matches {
		f, err := os.Open(n)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		r := bufio.NewReader(f)
		readData(suite, test, r)
	}
}
