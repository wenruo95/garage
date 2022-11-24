package tcp

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
)

type Server struct {
	status int32
	config *ServerConfig

	listener *net.TCPListener
	cancel   context.CancelFunc
}

func NewServer(options ...ServerOption) *Server {
	svr := new(Server)
	svr.status = int32(ServerStatusInit)

	svr.config = newServerConfig()
	for _, option := range options {
		option(svr.config)
	}

	return svr
}

func (svr *Server) Serve() error {
	if svr == nil || svr.config == nil {
		return errors.New("init server error")
	}
	addr, err := net.ResolveTCPAddr("tcp", svr.config.addr)
	if err != nil {
		return errors.New("resolve error:" + err.Error())
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return errors.New("listen error:" + err.Error())
	}
	svr.listener = listener

	if !atomic.CompareAndSwapInt32(&svr.status, int32(ServerStatusInit), int32(ServerStatusServing)) {
		return errors.New("invalid status:" + svr.Status().String())
	}

	ctx, cancel := context.WithCancel(context.Background())
	svr.cancel = cancel
	for {
		select {
		case <-ctx.Done():
			return errors.New("context has canceled")

		default:
			conn, err := listener.AcceptTCP()
			if err != nil {
				return errors.New("accept connection error:" + err.Error())
			}
			if svr.Status() == ServerStatusServing {
				go newClientConn(svr, conn).serve(ctx)
			}
		}
	}

	return nil
}

func (server *Server) Status() ServerStatus {
	return ServerStatus(atomic.LoadInt32(&server.status))
}

func (svr *Server) Close() error {
	if atomic.CompareAndSwapInt32(&svr.status, int32(ServerStatusServing), int32(ServerStatusClosed)) {
		svr.cancel()
		return svr.listener.Close()
	}

	atomic.StoreInt32(&svr.status, int32(ServerStatusClosed))
	return nil
}
