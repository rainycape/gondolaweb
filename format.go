package main

import (
	"bytes"
	"gnd.la/html"
	"go/printer"
	"go/token"
	"regexp"
)

var (
	commentRe          = regexp.MustCompile("//[^\n]+")
	multilineCommentRe = regexp.MustCompile("/\\*.*?\\*/")
)

func FormatNode(fset *token.FileSet, node interface{}) (string, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return "", err
	}
	escaped := html.Escape(buf.String())
	return FormatComments(escaped), nil
}

func FormatComments(text string) string {
	repl := func(s string) string {
		return "<span class=\"comments\">" + s + "</span>"
	}
	return multilineCommentRe.ReplaceAllStringFunc(commentRe.ReplaceAllStringFunc(text, repl), repl)
}
