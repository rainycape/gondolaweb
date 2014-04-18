package main

import (
	"gnd.la/app"
)

func MainHandler(ctx *app.Context) {
	ctx.WriteString("Hello world!")
}
