package packer

import (
	"image"

	"golang.org/x/image/draw"
)

func fastDraw(dst *image.NRGBA, r image.Rectangle, src image.Image) {
	w, h := r.Dx(), r.Dy()
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.BiLinear.Scale(img, image.Rect(0, 0, w, h), src, src.Bounds(), draw.Src, nil)
	drawCopySrc(dst, r, img, image.ZP)
}

func drawCopySrc(dst *image.NRGBA, r image.Rectangle, src *image.NRGBA, sp image.Point) {
	n, dy := 4*r.Dx(), r.Dy()
	d0 := dst.PixOffset(r.Min.X, r.Min.Y)
	s0 := src.PixOffset(sp.X, sp.Y)
	var ddelta, sdelta int
	if r.Min.Y <= sp.Y {
		ddelta = dst.Stride
		sdelta = src.Stride
	} else {
		// If the source start point is higher than the destination start
		// point, then we compose the rows in bottom-up order instead of
		// top-down. Unlike the drawCopyOver function, we don't have to check
		// the x coordinates because the built-in copy function can handle
		// overlapping slices.
		d0 += (dy - 1) * dst.Stride
		s0 += (dy - 1) * src.Stride
		ddelta = -dst.Stride
		sdelta = -src.Stride
	}
	for ; dy > 0; dy-- {
		copy(dst.Pix[d0:d0+n], src.Pix[s0:s0+n])
		d0 += ddelta
		s0 += sdelta
	}
}
