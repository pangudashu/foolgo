package foolgo

import (
	"net"
)

type FoolConn struct {
	net.Conn
}

func (this FoolConn) Close() error {
	connWg.Done()
	return this.Conn.Close()
}
