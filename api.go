package main

import (
	"gnd.la/app"
	"gnd.la/internal/project"
)

const (
	templateDownloadHandlerName = "template-download"
)

var (
	templates []*project.Template
)

func templateListHandler(ctx *app.Context) (interface{}, error) {
	return templates, nil
}

func templateDownloadHandler(ctx *app.Context) {
	name := ctx.IndexValue(0)
	for _, v := range templates {
		if v.Name == name {
			data, err := v.Data()
			if err != nil {
				panic(err)
			}
			ctx.Header().Set("Content-Type", "application/x-gzip")
			ctx.Write(data)
			return
		}
	}
	ctx.NotFound("template not found")
}
