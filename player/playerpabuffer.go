package player

import (
	"errors"
	"github.com/bobertlo/go-mpg123/mpg123"
	"sync"
)

// intermediate memory buffer reduces file io calls in audio callback
// this buffer will be held in memory while the file is playing, size accordingly
// the tradeoff with smaller buffer is more frequent file io

type mpgBuff struct {
	eof bool
	// max capacity: len(bufA) + len(bufB)
	size    int
	fileDec *mpg123.Decoder
	// two buffers, each half sized
	bufA *concBuff
	bufB *concBuff
	// boolean used to switch between buffers
	cur bool
}

type concBuff struct {
	sync.Mutex
	data []byte
	len  int
	pos  int
}

func (c *concBuff) next() byte {
	c.pos++
	return c.data[c.pos-1]
}

func (c *concBuff) has() bool {
	return c.pos < c.len
}

// mpgBuff implements io.Reader
func (b *mpgBuff) Read(p []byte) (int, error) {
	// copy bytes from memory buffer
	r := len(p)
	bf := b.buffer()
	bf.Lock()
	n := 0
	for ; n < r; n++ {
		if !bf.has() {
			// switch buffers
			bf.Unlock()
			bf = b.next()
			bf.Lock()
			if !bf.has() {
				bf.Unlock()
				// if second buffer is empty, we're done
				return n, errors.New("EOF")
			}
		}
		p[n] = bf.next()
	}
	bf.Unlock()
	return n, nil
}

// get intermediate buffer size
func (b *mpgBuff) BufferSize() int {
	return b.size
}

// set up and fill both buffers
func (b *mpgBuff) Prepare() {
	cap := b.BufferSize()
	// ensure capacity is divisible by 4
	cap += cap % 4
	// make sure we have an even cap value
	cap += (cap % 2)
	// get individual buffer length
	l := cap / 2
	b.eof = false
	// make and initialize buffers
	b.bufA = &concBuff{
		data: make([]byte, l)}
	b.fill(b.bufA)
	b.bufB = &concBuff{
		data: make([]byte, l)}
	b.fill(b.bufB)
}

// get current active reading buffer
func (b *mpgBuff) buffer() *concBuff {
	if !b.cur {
		return b.bufA
	}
	return b.bufB
}

// switch to next buffer, fill the other
func (b *mpgBuff) next() *concBuff {
	// request fill
	go func(buff *concBuff) {
		b.fill(buff)
	}(b.buffer())
	b.cur = !b.cur
	return b.buffer()
}

// fill buffer from mpg123 decoder without reallocating memory
func (b *mpgBuff) fill(buff *concBuff) {
	buff.Lock()
	n, err := b.fileDec.Read(buff.data)
	buff.len = n
	buff.pos = 0
	if err == mpg123.EOF {
		b.eof = true
	}
	buff.Unlock()
}
