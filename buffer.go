package main

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

type line struct {
	data []byte
	next *line
	prev *line
}

func (l *line) String() string {
	cur := "nil"
	next := "nil"
	prev := "nil"
	if l != nil {
		cur = string(l.data)
	}
	if l.next != nil {
		next = string(l.next.data)
	}
	if l.prev != nil {
		prev = string(l.prev.data)
	}
	return fmt.Sprintf("%s <- %s -> %s", prev, cur, next)
}

// Find a set of closest offsets for a given visual offset
func (l *line) find_closest_offsets(voffset int) (bo, co, vo int) {
	data := l.data
	for len(data) > 0 {
		var vodif int
		r, rlen := utf8.DecodeRune(data)
		data = data[rlen:]
		vodif = rune_advance_len(r, vo)
		if vo+vodif > voffset {
			return
		}

		bo += rlen
		co += 1
		vo += vodif
	}
	return
}

type buffer struct {
	first_line *line
	last_line  *line
	lines_n    int
	bytes_n    int

	words_cache       llrb_tree
	words_cache_valid bool
}

func new_empty_buffer() *buffer {
	b := new(buffer)
	l := new(line)
	l.next = nil
	l.prev = nil
	b.first_line = l
	b.last_line = l
	b.lines_n = 1
	return b
}

func (b *buffer) reader() *buffer_reader {
	return new_buffer_reader(b)
}

func (b *buffer) save_as(filename string) error {
	r := b.reader()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------
// buffer_reader
//----------------------------------------------------------------------------

type buffer_reader struct {
	buffer *buffer
	line   *line
	offset int
}

func new_buffer_reader(buffer *buffer) *buffer_reader {
	br := new(buffer_reader)
	br.buffer = buffer
	br.line = buffer.first_line
	br.offset = 0
	return br
}

func (br *buffer_reader) Read(data []byte) (int, error) {
	nread := 0
	for len(data) > 0 {
		if br.line == nil {
			return nread, io.EOF
		}

		// how much can we read from current line
		can_read := len(br.line.data) - br.offset
		if len(data) <= can_read {
			// if this is all we need, return
			n := copy(data, br.line.data[br.offset:])
			nread += n
			br.offset += n
			break
		}

		// otherwise try to read '\n' and jump to the next line
		n := copy(data, br.line.data[br.offset:])
		nread += n
		data = data[n:]
		if len(data) > 0 && br.line != br.buffer.last_line {
			data[0] = '\n'
			data = data[1:]
			nread++
		}

		br.line = br.line.next
		br.offset = 0
	}
	return nread, nil
}
