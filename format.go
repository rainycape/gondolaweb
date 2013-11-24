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

func commentRepl(text string) string {
	return "<span class=\"comments\">" + text + "</span>"
}

func FormatComments(text string) string {
	return multilineCommentRe.ReplaceAllStringFunc(commentRe.ReplaceAllStringFunc(text, commentRepl), commentRepl)
}
