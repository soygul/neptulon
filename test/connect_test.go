package test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/neptulon/client"
	"github.com/neptulon/neptulon/middleware"
)

// func TestConnectTCP(t *testing.T) {
// 	s := NewTCPServerHelper(t)
// 	defer s.Close()
// 	c := s.GetTCPClient()
// 	defer c.Close()
// }

func TestConnectTLS(t *testing.T) {
	sh := NewTLSServerHelper(t)
	defer sh.Close()

	sh.Server.MiddlewareIn(middleware.Echo)

	ch := sh.GetTLSClient(true)
	defer ch.Close()

	// todo: enable debug mode both on client & server if debug env var is defined during test launch

	var wg sync.WaitGroup
	msg := []byte("test message")

	ch.Client.MiddlewareIn(func(ctx *client.Ctx) {
		defer wg.Done()
		if !reflect.DeepEqual(ctx.Msg, msg) {
			t.Fatalf("expected: '%s', got: '%s'", msg, ctx.Msg)
		}
		ctx.Next()
	})

	wg.Add(1)
	ch.Send(msg)
	wg.Wait()
}
