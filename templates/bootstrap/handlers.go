package main

import (
	"net/http"

	"gnd.la/app"
)

func MainHandler(ctx *app.Context) {
	ctx.MustExecute("main.html", nil)
}

func ErrorHandler(ctx *app.Context, msg string, code int) bool {
	// Only send the 404 page if the error is a 404
	if code == http.StatusNotFound {
		ctx.MustExecute("404.html", nil)
		return true
	}
	// Otherwise, let the Gondola error processing take care of it
	return false
}
