package tcp

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type ClientConn struct {
	Conn    *net.TCPConn
	status  int32
	handler ClientHandler
	timeout time.Duration
	logger  *logrus.Entry
	cancel  context.CancelFunc
}

func newClientConn(server *Server, conn *net.TCPConn) *ClientConn {
	cc := new(ClientConn)
	cc.Conn = conn
	cc.status = int32(ServerConnStatusInit)
	cc.handler = server.config.handler

	if v := server.config.clientHBTimeout.Milliseconds(); v > MinClientHBTimeout {
		cc.timeout = time.Duration(v) * time.Millisecond
	} else {
		cc.timeout = time.Duration(MinClientHBTimeout) * time.Millisecond
	}

	cc.logger = logrus.WithField("client_addr", cc.Conn.RemoteAddr().String()).
		WithField("local_addr", cc.Conn.LocalAddr().String())

	return cc
}

func (cc *ClientConn) serve(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	cc.cancel = cancel
	cc.logger = cc.logger.WithContext(ctx)

	if err := cc.handler.OnConnect(cc); err != nil {
		cc.logger.Errorf("on_connect error:%v", err)
		return
	}

	timer := time.NewTimer(cc.timeout)
	for {
		select {
		case <-ctx.Done():
			cc.disconnect("context canceled")
			return

		case <-timer.C:
			cc.disconnect("conn heartbeat exceed " + cc.timeout.String())
			return

		default:
			body, messageFlag, txid, err := Unpack(cc.Conn)
			if err != nil {
				reason := "unpack " + cc.Conn.RemoteAddr().String() + " msg error:" + err.Error()
				cc.disconnect(reason)
				return
			}

			timer.Reset(cc.timeout)
			if err := cc.handler.OnMessage(cc, body, messageFlag, txid); err != nil {
				cc.disconnect(err.Error())
				cc.logger.Errorf("on_message error:%v", err)
				return
			}

		}
	}

}

func (cc *ClientConn) Status() ServerConnStatus {
	return ServerConnStatus(atomic.LoadInt32(&cc.status))
}

func (cc *ClientConn) Send(body []byte, messageFlag byte, txid uint32) error {
	if status := cc.Status(); status != ServerConnStatusServing {
		return errors.New("invalid conn statu:" + status.String())
	}
	return PackWrite(cc.Conn, messageFlag, txid, body)
}

func (cc *ClientConn) disconnect(reason string) {
	if err := cc.Close(); err != nil {
		cc.logger.Errorf("close error:%v", err)
	}
	cc.handler.OnDisconnect(cc, reason)
}

func (cc *ClientConn) Close() error {
	if !atomic.CompareAndSwapInt32(&cc.status, int32(ServerConnStatusServing), int32(ServerConnStatusClosed)) {
		return errors.New("client has stoped")
	}

	cc.cancel()
	return cc.Conn.Close()
}
