package packing

import (
	"math"
)

// ByArea implements sort Interface for []Block
// based on the Area of each block.
type ByArea []Block

func (a ByArea) Len() int      { return len(a) }
func (a ByArea) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool {
	iw, ih := a[i].Size()
	jw, jh := a[j].Size()
	return iw*ih > jw*jh
}

// ByMaxSide implements sort interface for []Block
// by comparing the maximum side (width or height)
// of each block
type ByMaxSide []Block

func (a ByMaxSide) Len() int      { return len(a) }
func (a ByMaxSide) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByMaxSide) Less(i, j int) bool {
	wi, hi := a[i].Size()
	wj, hj := a[j].Size()
	return math.Max(float64(wi), float64(hi)) > math.Max(float64(wj), float64(hj))
}
