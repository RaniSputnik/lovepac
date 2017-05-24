package packing

import "errors"

// ErrInputTooLarge means that a given block was larger than
// the max size of the packer - it can not possibly fit
var ErrInputTooLarge = errors.New("input too large, can not be processed")

// ErrOutOfRoom indicates a partial success, not all of the
// blocks fit into the maximum size of the packer
var ErrOutOfRoom = errors.New("out of room, could not place all blocks")

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

// Packer is the interface that wraps the Fit method.
//
// Fit sets a width and height for a packer then attempts
// to fit all blocks into the given region. If the blocks
// do not fit into the given area then an error is returned.
//
// Implementations should not retain blocks.
type Packer interface {
	Fit(width int, height int, blocks ...Block) error
}

type node struct {
	x, y int
	w, h int

	used  bool
	right *node
	down  *node
}
