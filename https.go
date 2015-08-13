package foolgo

import (
	"crypto/tls"
	"net"
)

func NewTlsListener(ln net.Listener, cert string, cert_key string) (tls_l net.Listener, err error) {
	config := &tls.Config{}
	config.NextProtos = []string{"http/1.1"}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(cert, cert_key)

	if err != nil {
		return nil, err
	}

	tls_l = tls.NewListener(ln, config)
	return tls_l, nil
}
