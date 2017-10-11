package config

const (
	// MaxClientPoolSize is the max size of client connection pool
	MaxClientPoolSize = 8

	// ClientReadTimeout is the timeout value for reading from connection
	ClientReadTimeout = 5 // seconds

	// ProxyAddr is the proxy listening port
	ProxyAddr = "127.0.0.1:8848"
)
