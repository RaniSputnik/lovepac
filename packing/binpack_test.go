package packing_test

import (
	"testing"

	. "github.com/RaniSputnik/lovepac/packing"
)

func TestBinPackingReturnsResults(t *testing.T) {
	blocks := []Block{
		&TestBlock{id: "1.png", w: 200, h: 200},
		&TestBlock{id: "2.png", w: 100, h: 100},
		&TestBlock{id: "3.png", w: 100, h: 50},
	}

	packer := NewBinPacker(300, 300)
	for _, block := range blocks {
		if err := packer.Pack(block); err != nil {
			t.Errorf("Expected that packer.Pack would not return an error but got %s", err.Error())
		}
	}

	for _, block := range blocks {
		testBlock := block.(*TestBlock)
		if testing.Verbose() {
			t.Logf("Testing block (%s), it has result: {%v,%v}", testBlock.id, testBlock.x, testBlock.y)
		}
		if !testBlock.placeWasCalled {
			t.Errorf("Block (%s) did not receive a result node", testBlock.id)
		}
	}
}

func TestBinPackingReturnsErrorIfInputBlockWillNeverFit(t *testing.T) {
	packer := NewBinPacker(100, 100)
	err := packer.Pack(&TestBlock{id: "doesnotfit.png", w: 200, h: 200})

	expected := ErrInputTooLarge
	if err != expected {
		t.Errorf("Expected packer.Pack to return '%v' but got '%v'", expected, err)
	}
}

func TestBinPackingReturnsErrorIfItRunsOutOfSpace(t *testing.T) {
	blocks := []Block{
		&TestBlock{id: "1.png", w: 200, h: 200},
		&TestBlock{id: "2.png", w: 100, h: 100},
		&TestBlock{id: "3.png", w: 100, h: 50},
	}

	packer := NewBinPacker(200, 200)
	err1 := packer.Pack(blocks[0])
	err2 := packer.Pack(blocks[1])

	if err1 != nil {
		t.Errorf("Expected packer.Pack of '1.png' to fit but got '%v'", err1)
	}

	if err2 != ErrOutOfRoom {
		t.Errorf("Expected packer.Pack of '2.png' to return '%v' but got '%v'", ErrOutOfRoom, err2)
	}
}

func TestBinPackingStillContinuesWhenRunOutOfSpace(t *testing.T) {
	blocks := map[Block]error{
		&TestBlock{id: "1.png", w: 200, h: 200}: nil,
		&TestBlock{id: "2.png", w: 200, h: 200}: ErrOutOfRoom,
		&TestBlock{id: "3.png", w: 100, h: 50}:  nil,
	}

	packer := NewBinPacker(300, 300)
	for block, expectedErr := range blocks {
		if err := packer.Pack(block); err != expectedErr {
			t.Errorf("Expected packer.Pack of block '%s' to return '%v' but got '%v'",
				block.(*TestBlock).id, expectedErr, err)
		}
	}

	for block, expectedErr := range blocks {
		testBlock := block.(*TestBlock)
		expectedToBePlaced := expectedErr == nil
		if testBlock.placeWasCalled != expectedToBePlaced {
			t.Errorf("Expected block (%s) placed to be '%t', but got '%t'",
				testBlock.id, expectedToBePlaced, testBlock.placeWasCalled)
		}
	}
}
