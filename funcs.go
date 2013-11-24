package main

import (
	"gnd.la/template"
	"go/doc"
	"strings"
)

func trim(s string, t string) string {
	return strings.Trim(s, t)
}

func funcId(fn *doc.Func) string {
	if fn.Recv != "" {
		recv := fn.Recv
		if recv[0] == '*' {
			recv = recv[1:]
		}
		return "type-" + recv + "-method-" + fn.Name
	}
	return "func-" + fn.Name
}

func typeId(fn *doc.Type) string {
	return "type-" + fn.Name
}

func init() {
	template.AddFuncs(template.FuncMap{
		"trim":    trim,
		"func_id": funcId,
		"type_id": typeId,
	})
}
