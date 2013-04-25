package main

import "testing"
import "strings"
import "io/ioutil"

func TestWrap(t *testing.T) {
	b := make([]byte, 1000)
	for i:=0 ; i<1000; i++ {
		b[i] = 'c'
	}
	s := strings.NewReader(string(b))
	wr, err := NewReader(s, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	s2, err := ioutil.ReadAll(wr)
	if err != nil {
		t.Fatalf(err.Error())
	}
	f := strings.Split(string(s2), "\n")
	if len(f[0]) != 10 {
		t.Errorf("should have been 10 long instead got %s", f[0])
	}

}
