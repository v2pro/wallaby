package config

const (
	// MaxClientPoolSize is the max size of client connection pool
	MaxClientPoolSize = 8

	// ClientReadTimeout is the timeout value for reading from connection
	ClientReadTimeout = 5 // seconds

	// ProxyAddr is the proxy listening ip:port
	ProxyAddr = "127.0.0.1:8868"

	// VersionHandlerAddr is the version handler listening ip:port
	VersionHandlerAddr = "127.0.0.1:8869"

	// ProxyServiceName is the name of service behind the proxy
	ProxyServiceName = "echo"

	// ProxyServiceConfig is the file path of Service Version Config
	ProxyServiceVersionConfig = "echo.json"

	// ProxyServiceConfig is the file path of Service Version Config
	ProxyBuildTimestamp = 1509598098
)
