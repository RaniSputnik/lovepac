package packing

import (
	"math"
)

// ByArea implements sort Interface for []Block
// based on the Area of each block.
type ByArea []Block

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool { return a[i].Width()*a[i].Height() > a[j].Width()*a[j].Height() }

// ByMaxSide implements sort interface for []Block
// by comparing the maximum side (width or height)
// of each block
type ByMaxSide []Block

func (a ByMaxSide) Len() int      { return len(a) }
func (a ByMaxSide) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByMaxSide) Less(i, j int) bool {
	wi := float64(a[i].Width())
	hi := float64(a[i].Height())
	wj := float64(a[j].Width())
	hj := float64(a[j].Height())
	return math.Max(wi, hi) > math.Max(wj, hj)
}
