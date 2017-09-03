package packing

import "errors"

// ErrInputTooLarge means that a given block was larger than
// the max size of the packer - it can not possibly fit
var ErrInputTooLarge = errors.New("input too large, can not be processed")

// ErrOutOfRoom indicates that the packer does not have enough
// room left for the block to be packed successfully.
var ErrOutOfRoom = errors.New("out of room")

// Block is the interface that represents a unit of space
// that can be packed alongside other Blocks.
//
// Size returns the width and height of the block.
//
// Place is called by the packer to indicate that the block
// has successfully been placed at the given position.
type Block interface {
	Size() (w int, h int)
	Place(x int, y int)
}

// Packer is the interface that wraps the Pack method.
type Packer interface {
	Pack(block Block) error
}

type node struct {
	x, y int
	w, h int

	used  bool
	right *node
	down  *node
}
