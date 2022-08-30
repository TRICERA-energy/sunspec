package sunspec

// Config is the configuration for a client or server.
type Config struct {
	// Endpoint specifics the sunspec host and is mandatory.
	// Currently only modbus tcp-networking is supported.
	// The schema must be host:port
	Endpoint string
}

// Client instantiates a new client from the given configuration.
func (o Config) Client() *Client {
	return &Client{client: newModbusClient(o.Endpoint)}
}

// Server instantiates a new server from the given configuration.
func (o Config) Server() *Server {
	return &Server{server: newModbusServer(o.Endpoint)}
}
