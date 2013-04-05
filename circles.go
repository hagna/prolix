package main

import "errors"

// CirclePoints prints all the symetrical points on the circle
// page 82 of Computer Graphics
func CirclePoints(x, y int, write_pixel func(x, y int)) {
	write_pixel(x, y)
	write_pixel(x, -y)
	write_pixel(-x, y)
	write_pixel(-x, -y)
	write_pixel(y, x)
	write_pixel(y, -x)
	write_pixel(-y, x)
	write_pixel(-y, -x)

}

func MidpointCircle(radius int, write_pixel func(x, y int)) {
	x := 0
	y := radius
	d := 5.0 / 4.0
	CirclePoints(x, y, write_pixel)

	for y > x {
		if d < 0 {
			d += 2.0*float64(x) + 3.0
		} else {
			d += 2.0*float64(x-y) + 5.0
			y--
		}
		x++
		CirclePoints(x, y, write_pixel)
	}
}

type xy struct {
	x, y int
}

type CircleVille struct {
	offset, r int
	bounds    map[int]xy //todo just an array here
}

func min(a []int) (res int) {
	res = a[0]
	for i := 1; i < len(a); i++ {
		if a[i] < res {
			res = a[i]
		}
	}
	return res
}

func max(a []int) (res int) {
	res = a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > res {
			res = a[i]
		}
	}
	return res
}

func makeCircleVille(radius int) CircleVille {
	m := make(map[int][]int)
	MidpointCircle(radius, func(x, y int) {
		m[y] = append(m[y], x)
	})
	bounds := make(map[int]xy)
	offset := 0
	for i := radius; i > -(radius + 1); i-- {

		for j := min(m[i]); j < max(m[i]); j++ {
			bounds[offset] = xy{j, i}
			offset++
		}
	}
	res := CircleVille{0, radius, bounds}
	return res
}

func (c *CircleVille) getxy(offset int) (x, y int, e error) {
	//todo overflow error
	x, y = 0, 0
	e = errors.New("overflow")
	if val, ok := c.bounds[offset]; ok {
		x, y = val.x, val.y
		e = nil
	}
	return
}
