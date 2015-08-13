package foolgo

import (
	"net"
	"os"
	"syscall"
)

type FoolListener struct {
	net.Listener
	stopped bool
}

func NewListener(addr string) (*FoolListener, error) {
	var l net.Listener
	var err error

	if isChild {
		if _, err := os.FindProcess(os.Getppid()); err != nil {
			return nil, err
		}

		f := os.NewFile(3, "")
		l, err = net.FileListener(f)
	} else {
		l, err = net.Listen("tcp", addr)
	}
	listener := &FoolListener{Listener: l, stopped: false}

	return listener, err
}

func (this *FoolListener) Accept() (c net.Conn, err error) {
	c, err = this.Listener.Accept()
	if err != nil {
		return c, err
	}

	c = FoolConn{
		Conn: c,
	}
	connWg.Add(1)
	return c, nil
}

func (this *FoolListener) Close() error {
	if this.stopped {
		return syscall.EINVAL
	}
	this.stopped = true

	return this.Listener.Close()
}

func (this *FoolListener) File() *os.File {
	// returns a dup(2) - FD_CLOEXEC flag *not* set
	l := this.Listener.(*net.TCPListener)
	file, _ := l.File()
	return file
}
