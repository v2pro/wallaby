package codec

import (
	"bufio"
	"errors"
	"io"
)

// Capture includes reader and the buffer
type Capture struct {
	reader *bufio.Reader
	bytes  []byte
}

// SetReader sets the reader and buffers
func (c *Capture) SetReader(reader io.Reader) error {
	if c.reader == nil {
		c.reader = bufio.NewReaderSize(&wrappedReader{
			capture: c,
			reader:  reader,
		}, 2048)
		c.bytes = make([]byte, 0, 2048)
	} else {
		if c.reader.Buffered() > 0 {
			return errors.New("has remaining buffer")
		}
		c.reader.Reset(&wrappedReader{
			capture: c,
			reader:  reader,
		})
		c.bytes = c.bytes[:]
	}
	return nil
}

// Reader get the Reader obj
func (c *Capture) Reader() *bufio.Reader {
	return c.reader
}

// Bytes get the buffer from reader
func (c *Capture) Bytes() []byte {
	captured := c.bytes[:len(c.bytes)-c.reader.Buffered()]
	c.bytes = c.bytes[len(c.bytes)-c.reader.Buffered():]
	return captured
}

// wrappedReader is a wrapper
type wrappedReader struct {
	capture *Capture
	reader  io.Reader
}

// Read reads from the input and append to wr
func (wr *wrappedReader) Read(p []byte) (n int, err error) {
	n, err = wr.reader.Read(p)
	wr.capture.bytes = append(wr.capture.bytes, p[:n]...)
	return
}
