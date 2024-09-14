package server

import (
	"crypto/tls"
	"fmt"
	"net"
)

//TODO Need to implement the Listener Interface from standard library in order to inject into the Serve methods

//TODO Look at channels for tracked listener how Serve interacts with listener and where channels can be used on
// on listener methods ~ example: select: <-Accept() ...? start time? <-Close() stop time...? would this handled
// by Serve?

// TODO Potentially implement local serve method on Conn

func TLSCert(config *tls.Config, certFile, keyFile string) (*tls.Config, error) {
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to configure certificates: %w", err)
	}
	return config, nil
}

type TracedListener struct {
	net.Listener
	Config    *tls.Config
	traceInfo string // <- will be custom data struct
}

func NewTracedListener(inner net.Listener, config *tls.Config) net.Listener {
	l := new(TracedListener)
	l.Listener = inner
	l.Config = config
	return l
}

// Accept Traced Accept Method -- returns a wrapped Server function which inits a conn and performs handshake
func (tl *TracedListener) Accept() (net.Conn, error) {
	c, err := tl.Listener.Accept()
	tl.traceInfo = c.LocalAddr().String()
	fmt.Println(tl.traceInfo)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Close Traced Close Method
func (tl *TracedListener) Close() error {
	tl.Listener.Close()
	return nil
}

// Custom serve function to include channels for metrics
