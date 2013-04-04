package main

import (
	"github.com/nsf/termbox-go"
	"time"
	//"strings"
	"log"
	"os"
)

func draw_n(cb []termbox.Cell, x, y, count int, color termbox.Attribute, ch rune) {
	for i := 0; i < count; i++ {
		termbox.SetCell(x+i, y, ch, color, termbox.ColorDefault)
	}
}

type struct Ev {
	termbox.EventKey
	ttype string
}

func keyb() chan string {
	res := make(chan string)
	go func() {
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}
			res <- ev
		}
		log.Println("done with switch")
		select {
		case <-c:
			log.Println("booom")
/*			n := len(strings.Fields(string(buf)))
			color := termbox.ColorWhite + termbox.Attribute(n/X)
			c := n / 10
			draw_n(cb, 0, Y-1, c, color, '#')
			termbox.Flush()*/
		default:
			log.Println("default")
		}

	}
	}()
	return res
}

func draw(s chan string) {
	x, y := 0, 10
	X, _ := termbox.Size()
	cb := termbox.CellBuffer()
	go func() {
		for {
			switch ev := <-s {
			case "insert":
				if x == X {
					copy(cb, cb[X:])
					x = 0
				} 
				termbox.SetCell(x, y, ev.Ch, termbox.ColorWhite, termbox.ColorDefault)
				x += 1
				termbox.Flush()

			}

		}
	}()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	logfile, err := os.Create("log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()


	buf := make([]rune, 1024*1024, 1024*1024)
	c := time.After(10 * time.Second)

}
