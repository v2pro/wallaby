package core

import "io"

type InboundRequest interface {
	Feature() map[string]string
	Write(writer io.Writer) error
}