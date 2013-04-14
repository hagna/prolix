package main

import (
	"log"
	"bufio"	
	"fmt"
	"os"
	"path/filepath"
	"io"
)

func convertFiles() {
	lst, err := filepath.Glob("a*")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("prolix.txt") //append if reasonable TODO
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	for i:=0; i<len(lst); i++ {
		fmt.Println(lst[i])
		st, err := os.Stat(lst[i])
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(f, fmt.Sprintf("%s\n", st.ModTime().String()))
		af, err := os.Open(lst[i])
		if err != nil {
			log.Println(err)
			continue
		}
		br := bufio.NewReader(af)
		for {
			s, err := br.ReadString('\n')
			fmt.Println(s)
			fmt.Fprintf(f, fmt.Sprintf("\n%s", s))
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				break
			}
		}
		af.Close()
	}
}
		
