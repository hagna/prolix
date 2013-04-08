package main

import "errors"

type RectVille struct {
	x, y, w, h int
}

func makeRectVille(x, y, w, h int) RectVille {
	return RectVille{x, y, w, h}
}

func (r *RectVille) getxy(offset int) (x, y int, e error) {
	d := offset / r.w
	y = r.y - d
	x = r.x + offset%r.w
	e = nil
	if y < (r.y - r.h) {
		e = errors.New("overflow")
	}
	return
}

// gives the new offset after a crlf at the 
// current offset
func (r *RectVille) crlf(oldoffset int) (offset int) {
	x, _, _ := r.getxy(oldoffset)
	d := x - r.x
	if d == 0 {
		d = oldoffset
	} else {
		d = oldoffset + (r.w - d)
	}
	return d
}
