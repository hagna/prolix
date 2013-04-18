package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"os"
	"unicode/utf8"
)

const (
	tabstop_length = 8
)

type view_location struct {
	cursor       cursor_location
	top_line     *line
	top_line_num int

	// Various cursor offsets from the beginning of the line:
	// 1. in characters
	// 2. in visual cells
	// An example would be the '\t' character, which gives 1 character
	// offset, but 'tabstop_length' visual cells offset.
	cursor_coffset int
	cursor_voffset int

	// This offset is different from these three above, because it's the
	// amount of visual cells you need to skip, before starting to show the
	// contents of the cursor line. The value stays as long as the cursor is
	// within the same line. When cursor jumps from one line to another, the
	// value is recalculated. The logic behind this variable is somewhat
	// close to the one behind the 'top_line' variable.
	line_voffset int

	// this one is used for choosing the best location while traversing
	// vertically, every time 'cursor_voffset' changes due to horizontal
	// movement, this one must be changed as well
	last_cursor_voffset int
}

type madman struct {
	view_location
	termbox_event chan termbox.Event
	buffer        *buffer
	quitflag      bool
	cursor        cursor_location
}

// When 'cursor_line' was changed, call this function to possibly adjust the
// 'top_line'.
func (m *madman) adjust_top_line() {
	/*	vt := v.vertical_threshold()
		top := v.top_line
		co := v.cursor.line_num - v.top_line_num
		h := v.height()

		if top.next != nil && co >= h-vt {
			v.move_top_line_n_times(co - (h - vt) + 1)
			v.dirty = dirty_everything
		}

		if top.prev != nil && co < vt {
			v.move_top_line_n_times(co - vt)
			v.dirty = dirty_everything
		}*/
}

// When 'cursor_voffset' was changed usually > 0, then call this function to
// possibly adjust 'line_voffset'.
func (m *madman) adjust_line_voffset() {
	/*	ht := v.horizontal_threshold()
		w := v.uibuf.Width
		vo := v.line_voffset
		cvo := v.cursor_voffset
		threshold := w - 1
		if vo != 0 {
			threshold = w - ht
		}

		if cvo-vo >= threshold {
			vo = cvo + (ht - w + 1)
		}

		if vo != 0 && cvo-vo < ht {
			vo = cvo - ht
			if vo < 0 {
				vo = 0
			}
		}

		if v.line_voffset != vo {
			v.line_voffset = vo
			v.dirty = dirty_everything
		}*/
}

func (m *madman) main_loop() {
	m.termbox_event = make(chan termbox.Event, 20)
	go func() {
		for {
			m.termbox_event <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-m.termbox_event:
			ok := m.handle_event(&ev)
			if !ok {
				return
			}
			m.consume_more_events()
			m.draw()
			termbox.Flush()
		}
	}
}

func (m *madman) consume_more_events() {
loop:
	for {
		select {
		case ev := <-m.termbox_event:
			ok := m.handle_event(&ev)
			if !ok {
				break loop
			}
		default:
			break loop
		}
	}
}

// Move cursor to the beginning of the file (buffer).
func (m *madman) move_cursor_beginning_of_file() {
	c := cursor_location{m.buffer.first_line, 1, 0}
	m.move_cursor_to(c)
}
func (m *madman) move_cursor_to(c cursor_location) {
	log.Println("move cursor to")
	//v.dirty |= dirty_status
	if c.boffset < 0 {
		bo, co, vo := c.line.find_closest_offsets(m.last_cursor_voffset)
		m.cursor.boffset = bo
		m.cursor_coffset = co
		m.cursor_voffset = vo
	} else {
		vo, co := c.voffset_coffset()
		m.cursor.boffset = c.boffset
		m.cursor_coffset = co
		m.cursor_voffset = vo
	}

	log.Println("after vo bo co")
	if c.boffset >= 0 {
		m.last_cursor_voffset = m.cursor_voffset
	}

	if c.line != m.cursor.line {
		if m.line_voffset != 0 {
			//v.dirty = dirty_everything
		}
		m.line_voffset = 0
	}
	m.cursor.line = c.line
	m.cursor.line_num = c.line_num
	m.adjust_line_voffset()
	m.adjust_top_line()
	log.Println("after adjust_top_line")
	log.Println("after adjust_top_line")

	/*if v.ac != nil {
		// update autocompletion on every cursor move
		ok := v.ac.update(v.cursor)
		if !ok {
			v.ac = nil
		}
	}*/
}

func (m *madman) insert_rune(r rune) {
	var data [utf8.UTFMax]byte
	l := utf8.EncodeRune(data[:], r)
	c := m.cursor
	log.Println(c)
	c.line.data = append(c.line.data, byte(r))
	c.boffset += l
	m.move_cursor_to(c)

}

func (m *madman) on_key(ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeyCtrlQ:
		m.quitflag = true
	case termbox.KeySpace:
		m.insert_rune(' ')
	}
	if ev.Ch != 0 {
		m.insert_rune(ev.Ch)
	}
}

func (m *madman) resize() {
}

func (m *madman) handle_event(ev *termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		m.on_key(ev)
		if m.quitflag {
			return false
		}
	case termbox.EventResize:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		m.resize()
	case termbox.EventError:
		log.Fatal(ev.Err)
	}

	return true
}
func (m *madman) cursor_position() (int, int) {
	y := m.cursor.line_num - m.top_line_num
	x := m.cursor_voffset - m.line_voffset
	return x, y
}
func (m *madman) draw() {
	cx, cy := m.cursor_position()
	termbox.SetCursor(cx, cy)
}

func new_madman() *madman {
	m := new(madman)
	m.buffer = new_empty_buffer()
	m.move_cursor_beginning_of_file()
	return m
}

func main() {
	logfile, err := os.Create("log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	if err = termbox.Init(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("termbox.Init()")
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)

	madman := new_madman()
	//madman.resize()
	madman.draw()
	termbox.SetCursor(madman.cursor_position())
	termbox.Flush()
	madman.main_loop()
}
