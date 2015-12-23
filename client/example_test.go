package client_test

import (
	"fmt"

	"github.com/neptulon/client"
)

// Example demonstrating the Neptulon client.
// Example assumes that there is a Neptulon server running on local network address 127.0.0.1:3001
// running a single echo middleware which echoes all incoming messages back.
func Example() {
	c := client.NewClient(nil, nil)
	c.MiddlewareIn(func(ctx *client.Ctx) {
		fmt.Println("Server's reply:", ctx.Msg)
		ctx.Next()
	})
	c.ConnectTCP("127.0.0.1:3001", false)
	c.Send([]byte("echo"))
	c.Close()
	// ** Output: Server's reply: echo
}
