package main

import (
	"image"
	"path"
	"strings"
)

// sprite implements the Block interface for packing
// and contains information about the image that it
// was constructed to represent
type sprite struct {
	path string
	img  image.Image
	x, y int
	w, h int
}

// Implement block interface
func (s *sprite) Width() int  { return s.w }
func (s *sprite) Height() int { return s.h }
func (s *sprite) Place(x int, y int) {
	s.x = x
	s.y = y
}

// Used for template rendering
func (s *sprite) Name() string { return strings.Replace(path.Base(s.path), path.Ext(s.path), "", 1) }
func (s *sprite) Left() int    { return s.x }
func (s *sprite) Top() int     { return s.y }
