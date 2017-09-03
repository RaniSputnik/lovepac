package packing_test

type TestBlock struct {
	id             string
	x, y           int
	w, h           int
	placeWasCalled bool
}

func (b *TestBlock) Size() (int, int) {
	return b.w, b.h
}

func (b *TestBlock) Place(x int, y int) {
	b.placeWasCalled = true
	b.x = x
	b.y = y
}
