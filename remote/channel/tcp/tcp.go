package tcp

import (
	"crypto/tls"
	"net"

	"github.com/yinqiwen/gsnova/common/logger"
	"github.com/yinqiwen/gsnova/common/mux"
	"github.com/yinqiwen/gsnova/remote"
	"github.com/yinqiwen/pmux"
)

func servTCP(lp net.Listener) {
	for {
		conn, err := lp.Accept()
		if nil != err {
			continue
		}
		session, err := pmux.Server(conn, remote.InitialPMuxConfig())
		if nil != err {
			logger.Error("[ERROR]Failed to create mux session for tcp server with reason:%v", err)
			continue
		}

		muxSession := &mux.ProxyMuxSession{Session: session}
		go remote.ServProxyMuxSession(muxSession)
	}
	//ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func StartTcpProxyServer(addr string) error {
	lp, err := net.Listen("tcp", addr)
	if nil != err {
		logger.Error("[ERROR]Failed to listen TCP address:%s with reason:%v", addr, err)
		return err
	}
	logger.Info("Listen on TCP address:%s", addr)
	servTCP(lp)
	return nil
}

func StartTLSProxyServer(addr string, config *tls.Config) error {
	lp, err := net.Listen("tcp", addr)
	if nil != err {
		logger.Error("[ERROR]Failed to listen TLS address:%s with reason:%v", addr, err)
		return err
	}
	lp = tls.NewListener(lp, config)
	logger.Info("Listen on TLS address:%s", addr)
	servTCP(lp)
	return nil
}
