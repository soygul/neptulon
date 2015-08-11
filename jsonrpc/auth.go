package jsonrpc

import "log"

// CertAuth is a TLS certificate authentication middleware for Neptulon JSON-RPC app.
type CertAuth struct {
}

// NewCertAuth creates and registers a new certificate authentication middleware instance with a Neptulon JSON-RPC app.
func NewCertAuth(app *App) (*CertAuth, error) {
	a := CertAuth{}
	app.ReqMiddleware(a.reqMiddleware)
	app.ResMiddleware(a.resMiddleware)
	app.NotMiddleware(a.notMiddleware)
	return &a, nil
}

func (a *CertAuth) reqMiddleware(ctx *ReqContext) {
	if _, ok := ctx.Conn.Session.Get("userid"); ok {
		return
	}

	// if provided, client certificate is verified by the TLS listener so the peerCerts list in the connection is trusted
	certs := ctx.Conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		log.Println("Invalid client-certificate authentication attempt:", ctx.Conn.RemoteAddr())
		ctx.Done = true
		ctx.Conn.Close()
		return
	}

	userID := certs[0].Subject.CommonName
	ctx.Conn.Session.Set("userid", userID)
	log.Println("Client-certificate authenticated:", ctx.Conn.RemoteAddr(), userID)
}

func (a *CertAuth) resMiddleware(ctx *ResContext) {
	if _, ok := ctx.Conn.Session.Get("userid"); ok {
		return
	}

	ctx.Done = true
	ctx.Conn.Close()
}

func (a *CertAuth) notMiddleware(ctx *NotContext) {
	if _, ok := ctx.Conn.Session.Get("userid"); ok {
		return
	}

	ctx.Done = true
	ctx.Conn.Close()
}
