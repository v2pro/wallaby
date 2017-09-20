package core

type InboundRequest interface {
	Feature() map[string]string
}