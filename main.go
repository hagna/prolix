package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var nomad = flag.Bool("notmad", false, "show's you aren't the madman")

func draw_n(x, y, count int, color termbox.Attribute, ch rune) {
	for i := 0; i < count; i++ {
		termbox.SetCell(x+i, y, ch, color, termbox.ColorDefault)
	}
}

func draw_words(X, Y, words int) {
	color := termbox.ColorDefault
	col := words / X
	rest := words - col*X

	draw_n(0, Y-1, rest, color, '#'+rune(col))
}

type Ev struct {
	termbox.Event
	ttype string
}

func keyb() chan Ev {
	res := make(chan Ev)
	var ev termbox.Event
	go func() {
	loop:
		for {
			switch ev = termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlQ {
					break loop
				}
				ttype := "insert"
				if ev.Key == termbox.KeyBackspace ||
					ev.Key == termbox.KeyBackspace2 {
					ttype = "backspace"
				}
				res <- Ev{ev, ttype}
			}

		}
		res <- Ev{ev, "quit"}
	}()
	return res
}

func writePixelFactory(offx, offy int) func(x, y int, ch rune) {
	X, Y := termbox.Size()
	ox := X/2 + offx
	oy := Y/2 + offy
	var res = func(x, y int, ch rune) {
		nx := ox + x
		ny := oy - y
		termbox.SetCell(nx, ny, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
	return res
}

func createAndSave(fname string, buf []rune) {
	outf, err := os.Create(fname)
	defer outf.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(outf, string(buf))
}

func saveBuf(fname string, buf []rune, modified bool) {
	if modified {
		createAndSave(fname, buf)
		//save it dang it
	}
}

func draw(s <-chan Ev) chan bool {

	X, Y := termbox.Size()
	offset := 0
	buf := make([]rune, 0, 0)
	res := make(chan bool)
	write_pixel := writePixelFactory(0, 0)
	cville := makeRectVille(-X/2+1, Y/2-2, X-3, Y-4)
	saveit := make(chan Ev)
	go func() {
		fname := nextFile()
		modified := true
		timer := time.After(10 * time.Second)
		unsaved := 0
		for {
			select {
			case ev := <-saveit:
				modified = true
				timer = time.After(10 * time.Second)
				unsaved++
				if ev.Key == termbox.KeySpace {
					buf = append(buf, ' ')
				} else if ev.Key == termbox.KeyBackspace ||
					ev.Key == termbox.KeyBackspace2 {
					buf = buf[:len(buf)-1]
				} else if ev.Ch > '!' && ev.Ch < '}' {
					buf = append(buf, ev.Ch)
				} else if ev.Key == termbox.KeyEnter {
					buf = append(buf, '\n')
				}
				if unsaved > 20 || ev.ttype == "quit" {
					saveBuf(fname, buf, modified)
					modified = false
					unsaved = 0
				}
			case <-timer:
				timer = time.After(10 * time.Second)
				saveBuf(fname, buf, modified)
				modified = false
			}
		}
	}()
	go func() {
		words := 0
		inword := true
	loop:
		for {
			select {
			case ev := <-s:
				saveit <- ev
				switch ev.ttype {
				case "insert":
					x, y, err := cville.getxy(offset)
					if err != nil {
						offset = 0
						x, y, err = cville.getxy(offset)
					}

					if ev.Key == termbox.KeySpace {
						if inword {
							words++
						}
						write_pixel(x, y, ' ')
						offset += 1
						inword = false
					} else if ev.Ch > '!' && ev.Ch < '}' {
						write_pixel(x, y, ev.Ch)
						offset += 1
						inword = true
					} else if ev.Key == termbox.KeyEnter {

						offset = cville.crlf(offset)

					}

					draw_words(X, Y, words)
					termbox.Flush()
				case "backspace":
					offset -= 1
					x, y, err := cville.getxy(offset)
					if err != nil {
						offset = 0
						x, y, err = cville.getxy(offset)
					}
					write_pixel(x, y, ' ')
					termbox.Flush()
				case "quit":
					break loop

				}
			}

		}
		res <- true
	}()
	return res
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

func main() {
	flag.Parse()
	if *nomad {
		convertFiles()
		return
	}
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
	keychan := keyb()
	<-draw(keychan)

}
