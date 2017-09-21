package codec

import (
	"io"
	"bufio"
	"errors"
)

type Capture struct {
	reader *bufio.Reader
	bytes  []byte
}

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

func (c *Capture) Reader() *bufio.Reader {
	return c.reader
}

func (c *Capture) Bytes() []byte {
	captured := c.bytes[:len(c.bytes)-c.reader.Buffered()]
	c.bytes = c.bytes[len(c.bytes)-c.reader.Buffered():]
	return captured
}

type wrappedReader struct {
	capture *Capture
	reader  io.Reader
}

func (wr *wrappedReader) Read(p []byte) (n int, err error) {
	n, err = wr.reader.Read(p)
	wr.capture.bytes = append(wr.capture.bytes, p[:n]...)
	return
}
