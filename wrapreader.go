package main

import (
	"io/ioutil"
	"io"
	"fmt"
	)

type wrapreader struct {
	data []byte
	wrapmark int
	offset int
}

func NewReader(r io.Reader, wrapmark int) (res *wrapreader, err error) {
	res = new(wrapreader)
	res.data, err = ioutil.ReadAll(r)
	res.wrapmark = wrapmark
	res.offset = 0
	return
}

func (wr *wrapreader) Read(data []byte) (n int, err error) {
	err = nil
	n = 0
	for len(data) > 0 {
		if wr.offset == len(wr.data) {
			err = io.EOF
			break
		}
		if wr.offset%wr.wrapmark == 0 {
			data[0] = '\n'
		} else {
			data[0] = wr.data[wr.offset]
			wr.offset++
		}
			n++
		data = data[1:]
	}
	err = io.EOF
	fmt.Printf("len of data is %d\n", len(data))
	fmt.Printf("len of wr.data is %d\n", len(wr.data))
	fmt.Println(err)
	fmt.Printf("read %d bytes\n", n)
	return
}
