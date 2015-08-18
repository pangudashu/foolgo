package foolgo

import (
	"fmt"
	"net"
)

type FoolConn struct {
	net.Conn
}

func (this FoolConn) Close() error {
	fmt.Println("close")
	connWg.Done()
	return this.Conn.Close()
}
