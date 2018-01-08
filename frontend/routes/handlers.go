package routes

import macaron "gopkg.in/macaron.v1"

func pingHandler(ctx *macaron.Context) {
	type Ping struct {
		Ping string
	}
	ctx.JSON(200, Ping{"Hello, World!"})
}
