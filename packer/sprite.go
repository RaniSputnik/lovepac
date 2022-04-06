package packer

import (
	"path"
	"strings"
)

// sprite implements the Block interface for packing
// and contains information about the image that it
// was constructed to represent
type sprite struct {
	Asset
	path    string
	x, y    int
	w, h    int
	padding int
	placed  bool
}

// Implement block interface
func (s *sprite) Size() (int, int) {
	return s.w + s.padding, s.h + s.padding
}
func (s *sprite) Place(x int, y int) {
	s.x = x + s.padding
	s.y = y + s.padding
	s.placed = true
}

// Used for template rendering
func (s *sprite) Name() string        { return strings.Replace(path.Base(s.path), path.Ext(s.path), "", 1) }
func (s *sprite) DisplayName() string { return strings.Replace(s.path, path.Ext(s.path), "", 1) }
func (s *sprite) Left() int           { return s.x }
func (s *sprite) Top() int            { return s.y }
func (s *sprite) Width() int          { return s.w }
func (s *sprite) Height() int         { return s.h }
