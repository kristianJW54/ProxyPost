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

//------------------ To Implement -----------------//

//func (srv *TracedServer) ListenAndServe() error {
//	if srv.shuttingDown() {
//		return ErrServerClosed
//	}
//	addr := srv.Addr
//	if addr == "" {
//		addr = ":http"
//	}
//	ln, err := net.Listen("tcp", addr)
//	if err != nil {
//		return err
//	}
//	return srv.Serve(ln) <- Will need to create in order to control client connection channels
//}

const bufferBeforeChunkingSize = 2048

//TODO Create response struct and chunk struct with a chunk writer

type TracedListener struct {
	net.Listener
	Config      *tls.Config
	TraceInfo   map[int]string // <- will be custom data struct
	ClientCount int
}

func NewTracedListener(inner net.Listener, config *tls.Config) net.Listener {
	l := new(TracedListener)
	l.TraceInfo = make(map[int]string)
	l.ClientCount = 0
	l.Listener = inner
	l.Config = config
	return l
}

// Accept Traced Accept Method -- returns a wrapped Server function which inits a conn and performs handshake
func (tl *TracedListener) Accept() (net.Conn, error) {
	c, err := tl.Listener.Accept()
	if err != nil {
		return nil, err
	}
	tl.ClientCount++
	tl.TraceInfo[tl.ClientCount] = c.RemoteAddr().String()
	fmt.Println(tl.TraceInfo[tl.ClientCount])
	return c, nil
}

// Close Traced Close Method
func (tl *TracedListener) Close() error {
	err := tl.Listener.Close()
	if err != nil {
		return err
	}
	return nil
}

// Custom serve function to include channels for metrics
