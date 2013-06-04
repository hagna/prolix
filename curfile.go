package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type curfile struct {
	files []string
	i	int
}

func newcurfile() *curfile {
	c := new(curfile)
}

func nextFile() string {
	lst, err := filepath.Glob("a*")
	if err != nil {
		log.Fatal(err)
	}
	a := []int{}
	for i := 0; i < len(lst); i++ {
		v, err := strconv.Atoi(lst[i][1:])
		if err != nil {
			log.Println(err)
			continue
		}
		a = append(a, v)
	}
	if len(a) > 0 {
		lastfile := fmt.Sprintf("a%d", max(a))
		finf, err := os.Stat(lastfile)
		if err != nil {
			log.Println(err)
		}
		log.Printf("size of last file was %d", finf.Size())
		if finf.Size() < 3 {
			log.Printf("recycling last file %s", lastfile)
			return lastfile
		}
		return fmt.Sprintf("a%d", max(a)+1)
	}
	return "a1"
}
