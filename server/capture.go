package server

import "io"

type capture struct {
	reader io.Reader
	bytesRead []byte
}

func (c *capture) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	c.bytesRead = append(c.bytesRead, p[:n]...)
	return
}


