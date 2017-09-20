package core

import "io"

type Packet interface {
	Feature() map[string]string
	Write(io.Writer) error
}