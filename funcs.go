package main

import (
	"gnd.la/template"
	"go/doc"
	"strings"
)

const (
	constPrefix = "const-"
	varPrefix   = "var-"
)

func ConstId(name string) string {
	return constPrefix + name
}

func VarId(name string) string {
	return varPrefix + name
}

func FuncId(name string) string {
	return "func-" + name
}

func TypeId(name string) string {
	return "type-" + name
}

func MethodId(typ string, name string) string {
	return "type-" + typ + "-method-" + name
}

func trim(s string, t string) string {
	return strings.Trim(s, t)
}

func funcId(fn *doc.Func) string {
	if fn.Recv != "" {
		recv := fn.Recv
		if recv[0] == '*' {
			recv = recv[1:]
		}
		return MethodId(recv, fn.Name)
	}
	return FuncId(fn.Name)
}

func typeId(typ *doc.Type) string {
	return TypeId(typ.Name)
}

func init() {
	template.AddFuncs(template.FuncMap{
		"trim":    trim,
		"func_id": funcId,
		"type_id": typeId,
	})
}
