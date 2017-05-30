package packing

type BinPacker struct {
	root *node
}

// NewBinPacker returns a packer with the given width and height
func NewBinPacker(width, height int) *BinPacker {
	return &BinPacker{
		root: &node{x: 0, y: 0, w: width, h: height},
	}
}

// Size returns the width and height of the BinPacker
func (b *BinPacker) Size() (int, int) { return b.root.w, b.root.h }

// Width returns the width of the BinPacker (immutable)
func (b *BinPacker) Width() int { return b.root.w }

// Height returns the height of the BinPacker (immutable)
func (b *BinPacker) Height() int { return b.root.h }

// Pack implements the Packer interface
func (b *BinPacker) Pack(block Block) error {
	bw, bh := block.Size()
	if bw > b.root.w || bh > b.root.h {
		return ErrInputTooLarge
	}

	if n := b.findNode(b.root, bw, bh); n != nil {
		b.splitNode(n, bw, bh)
		block.Place(n.x, n.y)
	} else {
		return ErrOutOfRoom
	}

	return nil
}

func (b *BinPacker) findNode(root *node, w int, h int) *node {
	if root.used {
		if r := b.findNode(root.right, w, h); r != nil {
			return r
		}
		return b.findNode(root.down, w, h)
	} else if (w <= root.w) && (h <= root.h) {
		return root
	} else {
		return nil
	}
}

func (b *BinPacker) splitNode(n *node, w int, h int) {
	n.used = true
	n.right = &node{x: n.x + w, y: n.y, w: n.w - w, h: h}
	n.down = &node{x: n.x, y: n.y + h, w: n.w, h: n.h - h}
}
