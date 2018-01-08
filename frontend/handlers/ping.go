package handlers

import (
	m "gopkg.in/macaron.v1"
)

// SendPing sends a simple JSON demonstarting the server is alive.
// This looks like: "{ Ping: Hello, World! }"
func SendPing(ctx *m.Context) {
	ctx.JSON(200, Ping{"Hello, World!"})
}
