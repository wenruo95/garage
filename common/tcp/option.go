package tcp

import "time"

type ServerConfig struct {
	addr    string
	handler ClientHandler

	clientHBTimeout time.Duration // client heartbeat timeout
}

func newServerConfig() *ServerConfig {
	return nil
}

type ClientHandler interface {
	OnConnect(conn *ClientConn) error
	OnMessage(conn *ClientConn, body []byte, messageFlag byte, txid uint32) error
	OnDisconnect(conn *ClientConn, reason string) error
}

type ServerOption func(server *ServerConfig)

func WithServerAddr(addr string) ServerOption {
	return func(server *ServerConfig) {
		server.addr = addr
	}
}

func WithServerHandler(handler ClientHandler) ServerOption {
	return func(server *ServerConfig) {
		server.handler = handler
	}
}

// client_hb_timeout default 3 * 60s
func WithClientHBTimeout(timeout time.Duration) ServerOption {
	return func(server *ServerConfig) {
		server.clientHBTimeout = timeout
	}
}
