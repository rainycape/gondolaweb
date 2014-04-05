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
	key := "proj-tmpl-" + name
	c := ctx.Cache()
	data, _ := c.GetBytes(key)
	if data == nil {
		for _, v := range templates {
			if v.Name == name {
				var err error
				data, err = v.Data()
				if err != nil {
					panic(err)
				}
				break
			}
		}
		if data == nil {
			ctx.NotFound("template not found")
		}
		c.SetBytes(key, data, 0)
	}
	ctx.Header().Set("Content-Type", "application/x-gzip")
	ctx.Write(data)
}
