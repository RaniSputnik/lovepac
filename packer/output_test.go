package packer_test

import (
	"bytes"
	"io"
	"sync"
)

type OutputRecorder struct {
	writers map[string]*bufferWithClose
	*sync.Mutex
}

type bufferWithClose struct {
	*bytes.Buffer
}

func (b *bufferWithClose) Close() error { return nil }

func (r *OutputRecorder) GetWriter(filename string) (io.WriteCloser, error) {
	buffer := &bufferWithClose{bytes.NewBufferString("")}
	r.Lock()
	r.writers[filename] = buffer
	r.Unlock()
	return buffer, nil
}

func (r *OutputRecorder) Got() map[string]*bytes.Buffer {
	r.Lock()
	results := map[string]*bytes.Buffer{}
	for key, val := range r.writers {
		results[key] = val.Buffer
	}
	r.Unlock()
	return results
}

func NewOutputRecorder() *OutputRecorder {
	return &OutputRecorder{map[string]*bufferWithClose{}, &sync.Mutex{}}
}
