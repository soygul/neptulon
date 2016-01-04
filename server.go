// Package neptulon is a RPC framework with middleware support.
package neptulon

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/neptulon/cmap"

	"golang.org/x/net/websocket"
)

// Server is a Neptulon server.
type Server struct {
	addr       string
	conns      *cmap.CMap // conn ID -> Conn
	middleware []func(ctx *ReqCtx) error
	wsConfig   websocket.Config
}

// NewServer creates a new Neptulon server.
func NewServer(addr string) *Server {
	return &Server{
		addr:  addr,
		conns: cmap.New(),
	}
}

// UseTLS enables Transport Layer Security for the connection.
// cert, key = Server certificate/private key pair.
// clientCACert = Optional certificate for verifying client certificates.
// All certificates/private keys are in PEM encoded X.509 format.
func (s *Server) UseTLS(cert, privKey, clientCACert []byte) error {
	tlsCert, err := tls.X509KeyPair(cert, privKey)
	if err != nil {
		return fmt.Errorf("failed to parse the server certificate or the private key: %v", err)
	}

	c, _ := pem.Decode(cert)
	if tlsCert.Leaf, err = x509.ParseCertificate(c.Bytes); err != nil {
		return fmt.Errorf("failed to parse the server certificate: %v", err)
	}

	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(clientCACert)
	if !ok {
		return fmt.Errorf("failed to parse the CA certificate: %v", err)
	}

	s.wsConfig.TlsConfig = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    pool,
		ClientAuth:   tls.VerifyClientCertIfGiven,
	}

	return nil
}

// Middleware registers middleware to handle incoming request messages.
func (s *Server) Middleware(middleware ...func(ctx *ReqCtx) error) {
	s.middleware = append(s.middleware, middleware...)
}

// Start the Neptulon server. This function blocks until server is closed.
func (s *Server) Start() error {
	http.Handle("/", websocket.Server{
		Config:  s.wsConfig,
		Handler: s.wsHandler,
		Handshake: func(config *websocket.Config, req *http.Request) error {
			config.Origin, _ = url.Parse(req.RemoteAddr) // we're interested in remote address and not origin header text
			return nil
		},
	})
	log.Println("Neptulon server started at address:", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) wsHandler(ws *websocket.Conn) {
	log.Println("Client connected:", ws.RemoteAddr())
	c, err := NewConn(ws, s.middleware)
	if err != nil {
		log.Println("Error while accepting connection:", err)
		return
	}

	s.conns.Set(c.ID, c)
	c.StartReceive()
	s.conns.Delete(c.ID)
	log.Println("Connection closed:", ws.RemoteAddr())
}

// Close closes the network listener and the active connections.
// func (s *Server) Close() error {
// 	err := s.listener.Close()
//
// 	// close all active connections discarding any read/writes that is going on currently
// 	// this is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server
// 	s.clients.Range(func(c interface{}) {
// 		c.(*Client).Close()
// 	})
//
// 	if err != nil {
// 		return fmt.Errorf("And error occured before or while stopping the server: %v", err)
// 	}
//
// 	return nil
// }
