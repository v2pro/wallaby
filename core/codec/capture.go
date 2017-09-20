package codec

import (
	"io"
	"bufio"
)

type Capture struct {
	reader *bufio.Reader
	bytes []byte
}

func (c *Capture) Reset(reader io.Reader) {
	if c.reader == nil {
		c.reader = bufio.NewReaderSize(&wrappedReader{
			capture: c,
			reader: reader,
		}, 2048)
		c.bytes = make([]byte, 0, 2048)
	} else {
		c.reader.Reset(&wrappedReader{
			capture: c,
			reader: reader,
		})
		c.bytes = c.bytes[:]
	}
}

func (c *Capture) Reader() *bufio.Reader {
	return c.reader
}

func (c *Capture) Bytes() []byte {
	return c.bytes
}

type wrappedReader struct {
	capture *Capture
	reader io.Reader
}

func (wr *wrappedReader) Read(p []byte) (n int, err error) {
	n, err = wr.reader.Read(p)
	wr.capture.bytes = append(wr.capture.bytes, p[:n]...)
	return
}