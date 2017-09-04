package player

import (
	"errors"
	"github.com/bobertlo/go-mpg123/mpg123"
	"sync"
)

// intermediate memory buffer reduces file io calls in audio callback
// this buffer will be held in memory while the file is playing, size accordingly
// the tradeoff with smaller buffer is more frequent file io
const interBuffSize = 10000

type mpgBuff struct {
	eof bool
	// max capacity: len(bufA) + len(bufB)
	cap     int
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
	n := 0
	for ; n < r; n++ {
		if !bf.has() {
			// switch buffers
			bf = b.next()
			if !bf.has() {
				// if second buffer is empty, we're done
				return 0, errors.New("EOF")
			}
		}
		p[n] = bf.next()
	}
	return n, nil
}

// set up and fill both buffers
func (b *mpgBuff) Prepare() {
	// make sure we have an even cap value
	b.cap = b.cap + (b.cap % 2)
	// get individual buffer length
	l := b.cap / 2
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
	// request fill, unless we're EOF
	if !b.eof {
		go func(buff *concBuff) {
			b.fill(buff)
		}(b.buffer())
	}
	b.cur = !b.cur
	return b.buffer()
}

// fill buffer from mpg123 decoder without reallocating memory
func (b *mpgBuff) fill(buff *concBuff) {
	buff.Lock()
	// reset pos
	buff.pos = 0
	n, err := b.fileDec.Read(buff.data)
	buff.len = n
	if err == mpg123.EOF {
		b.eof = true
	}
	buff.Unlock()
}
