package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		return
	}

	f, _ := os.Open(flag.Arg(0))
	datas, _ := ioutil.ReadAll(f)
	r := newMyReader(datas)
	// r := f
	p := make([]byte, 3)
	n, err := r.Read(p)
	fmt.Println(n, err, string(p))
	n, err = r.Read(p)
	fmt.Println(n, err, string(p))

	// r.Seek(2, 0)

	n, err = r.Read(p)
	fmt.Println(n, err, string(p[:n]))
	n, err = r.Read(p)
	fmt.Println(n, err, string(p[:n]))
	n, err = r.Read(p)
	fmt.Println(n, err, string(p[:n]))
}

type tuple struct {
	data []byte
	n    int
	err  error
}

type myReader struct {
	b chan tuple
	// q      chan bool
	datas  []byte
	offset int64
}

func newMyReader(datas []byte) io.ReadSeeker {
	r := new(myReader)
	r.datas = datas
	return r
}

func (r *myReader) start(p []byte) {
	r.b = make(chan tuple)
	d := make([]byte, len(p))
	go func() {
		for {
			for r.offset < int64(len(r.datas)) {

				// fmt.Println("read at", r.offset)
				n := int64(len(p))
				var err error = nil
				if n+r.offset > int64(len(r.datas)) {
					n = int64(len(r.datas)) - r.offset
					err = io.EOF
				}
				copy(p, r.datas[r.offset:r.offset+n])
				r.offset += n
				r.b <- tuple{d, int(n), err}
			}

			r.b <- tuple{[]byte{}, 0, io.EOF}
		}
	}()
}

func (r *myReader) Read(p []byte) (int, error) {
	r.start(p)

	s, _ := <-r.b

	if s.err != nil {
		if s.err == io.EOF {
			return s.n, s.err
		}
		return 0, s.err
	}
	return s.n, nil
}

func (r *myReader) Seek(offset int64, whence int) (int64, error) {
	// close(r.b)
	r.offset = offset
	// fmt.Println("set to", r.offset)
	// r.start()
	return offset, nil
}
