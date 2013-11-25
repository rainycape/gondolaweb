package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gnd.la/mux"
	"gnd.la/util"
	"gnd.la/util/astutil"
	"gnd.la/util/pkgutil"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	srcDir    = util.RelativePath("src")
	cwd       = util.RelativePath("..")
	srcPrefix = srcDir[len(cwd)+1:]
	valueRe   = regexp.MustCompile("([A-Z]\\w+)((?:\\s+<a.*?</a>)?\\s+=)")
)

func noBuildable(err error) bool {
	return strings.Contains(err.Error(), "no buildable")
}

func buildContext() *build.Context {
	ctx := build.Default
	ctx.GOPATH = util.RelativePath(".")
	return &ctx
}

type Package struct {
	ctx      *mux.Context
	fset     *token.FileSet
	name     string
	bpkg     *build.Package
	apkg     *ast.Package
	dpkg     *doc.Package
	Packages []*Package
}

func (p *Package) symbolHref(symbol string) string {
	key := symbol
	if key[len(key)-1] == ')' && key[len(key)-2] == '(' {
		key = key[:len(key)-2]
	}
	if obj := p.apkg.Scope.Objects[key]; obj != nil {
		switch obj.Kind {
		case ast.Typ:
			return "#type-" + key
		case ast.Fun:
			return "#func-" + key
		case ast.Con:
			return "#const-" + key
		case ast.Var:
			return "#var-" + key
		}
	}
	if dot := strings.IndexByte(key, '.'); dot > 0 {
		tn := key[:dot]
		fn := key[dot+1:]
		if obj := p.apkg.Scope.Objects[tn]; obj != nil && obj.Kind == ast.Typ {
			for _, v := range p.dpkg.Types {
				if v.Name == tn {
					for _, m := range v.Methods {
						if m.Name == fn {
							return "#type-" + tn + "-method-" + fn
						}
					}
					return ""
				}
			}
		}
	}
	return ""
}

func (p *Package) href(word string, scope string) string {
	slash := strings.IndexByte(word, '/')
	dot := strings.IndexByte(word, '.')
	if slash > 0 || dot > 0 {
		// Check if there's a type or function mentioned
		// after the package.
		if pn, tn := pkgutil.SplitQualifiedName(word); pn != "" && tn != "" {
			if pn[0] == '*' {
				pn = pn[1:]
			}
			if pkg, err := ImportPackage(p.ctx, pn); err == nil {
				if sr := pkg.symbolHref(tn); sr != "" {
					return p.ctx.MustReverse("doc", pn) + sr
				}
			}
			if pn == p.dpkg.Name {
				return p.symbolHref(tn)
			}
		} else if _, err := buildContext().Import(word, "", build.FindOnly); err == nil {
			return p.ctx.MustReverse("doc", word)
		}
	}
	if dot > 0 {
		// Check the package imports, to see if any of them matches
		// TODO: Check for packages imported with a different local
		// name.
		base := word[:dot]
		for _, v := range p.bpkg.Imports {
			if path.Base(v) == base && v != base {
				return p.href(v+"."+word[dot+1:], scope)
			}
		}
	}
	if word[0]&0x20 == 0 {
		// Uppercase
		if scope != "" {
			if href := p.symbolHref(scope + "." + word); href != "" {
				return href
			}
		}
		return p.symbolHref(word)
	}
	return ""
}

func (p *Package) writeWord(bw *bufio.Writer, buf *bytes.Buffer, scope string, ignored map[string]struct{}) {
	if word := buf.String(); word != "" {
		if _, ign := ignored[word]; ign {
			bw.WriteString(word)
		} else {
			if href := p.href(word, scope); href != "" {
				bw.WriteString("<a href=\"")
				bw.WriteString(href)
				bw.WriteString("\">")
				bw.WriteString(word)
				bw.WriteString("</a>")
			} else {
				bw.WriteString(word)
			}
		}
	}
}

func (p *Package) linkify(w io.Writer, input string, scope string, ignored map[string]struct{}) error {
	bw := bufio.NewWriterSize(w, 512)
	var buf bytes.Buffer
	for ii := 0; ii < len(input); ii++ {
		c := input[ii]
		switch c {
		// Include * and & in the list of stop characters,
		// so pointers get the link for the pointed type.
		// Include ;, so escaped amperstands do not end up
		// in the type names.
		case ',', ' ', '\n', '\t', '(', ')', '*', '&', '{', '}', ';':
			p.writeWord(bw, &buf, scope, ignored)
			bw.WriteByte(c)
			buf.Reset()
		case '.':
			if next := ii + 1; next < len(input) {
				if nc := input[next]; nc == ' ' || nc == '\t' || nc == '\n' {
					p.writeWord(bw, &buf, scope, ignored)
					bw.WriteByte(c)
					buf.Reset()
					continue
				}
			}
			fallthrough
		default:
			buf.WriteByte(c)
		}
	}
	p.writeWord(bw, &buf, scope, ignored)
	return bw.Flush()
}

func (p *Package) Name() string {
	if p.bpkg != nil {
		return p.bpkg.Name
	}
	return p.name
}

func (p *Package) ImportPath() string {
	if p.bpkg != nil {
		path := p.bpkg.ImportPath
		if strings.HasPrefix(path, srcPrefix) {
			path = path[len(srcPrefix)+1:]
		}
		return path
	}
	return ""
}

func (p *Package) Synopsis() string {
	if p.dpkg != nil {
		return doc.Synopsis(p.dpkg.Doc)
	}
	return ""
}

func (p *Package) Filenames() []string {
	if b := p.bpkg; b != nil {
		var files []string
		for _, v := range [][]string{b.GoFiles, b.CgoFiles, b.IgnoredGoFiles, b.CFiles, b.CXXFiles, b.HFiles, b.SFiles, b.SysoFiles, b.SwigFiles, b.SwigCXXFiles} {
			files = append(files, v...)
		}
		sort.Strings(files)
		return files
	}
	return nil
}

func (p *Package) Position(n ast.Node) string {
	pos := p.fset.Position(n.Pos())
	filename := pos.Filename
	if strings.HasPrefix(filename, srcDir) {
		filename = filename[len(srcDir)+1:]
	}
	if gr := buildContext().GOROOT; strings.HasPrefix(filename, gr) {
		// Skip the src/ after GOROOT
		filename = "go" + filename[len(gr)+4:]
	}
	return fmt.Sprintf("%s#line-%d", filename, pos.Line)
}

func (p *Package) HasDoc() bool {
	return p.dpkg != nil && strings.TrimSpace(doc.Synopsis(p.dpkg.Doc)) != strings.TrimSpace(p.dpkg.Doc)
}

func (p *Package) Doc() *doc.Package {
	return p.dpkg
}

func (p *Package) html(text string, scope string, ignored map[string]struct{}) template.HTML {
	var buf bytes.Buffer
	doc.ToHTML(&buf, text, nil)
	var out bytes.Buffer
	p.linkify(&out, buf.String(), scope, ignored)
	return template.HTML(out.String())
}

func (p *Package) HTML(text string) template.HTML {
	return p.html(text, "", nil)
}

func (p *Package) scopeParameters(node interface{}) (string, map[string]struct{}) {
	var scope string
	var ignored map[string]struct{}
	ignore := func(x string) {
		if ignored == nil {
			ignored = make(map[string]struct{})
		}
		ignored[x] = struct{}{}
	}
	switch n := node.(type) {
	case *ast.FuncDecl:
		ignore(n.Name.Name)
		if n.Recv != nil {
			scope = astutil.Ident(n.Recv.List[0].Type)
			if scope != "" && scope[0] == '*' {
				scope = scope[1:]
			}
		}
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				scope = s.Name.Name
				ignore(scope)
			case *ast.ValueSpec:
				for _, name := range s.Names {
					ignore(name.Name)
				}
			}
		}
	}
	return scope, ignored
}

func (p *Package) ScopedHTML(text string, scope interface{}) template.HTML {
	name, ignored := p.scopeParameters(scope)
	return p.html(text, name, ignored)
}

func (p *Package) HTMLDoc() template.HTML {
	return p.HTML(p.dpkg.Doc)
}

func (p *Package) HTMLDecl(node interface{}) (template.HTML, error) {
	s, err := FormatNode(p.fset, node)
	if err == nil {
		name, ignored := p.scopeParameters(node)
		var buf bytes.Buffer
		p.linkify(&buf, s, name, ignored)
		s = buf.String()
	}
	if strings.HasPrefix(s, "const ") {
		s = valueRe.ReplaceAllString(s, "<span id=\"const-${1}\">${1}</span>${2}")
	} else if strings.HasPrefix(s, "var ") {
		s = valueRe.ReplaceAllString(s, "<span id=\"var-${1}\">${1}</span>${2}")
	}
	return template.HTML(s), err
}

func pkgImporter(imports map[string]*ast.Object, path string) (*ast.Object, error) {
	pkg := imports[path]
	if pkg == nil {
		// note that strings.LastIndex returns -1 if there is no "/"
		pkg = ast.NewObj(ast.Pkg, path[strings.LastIndex(path, "/")+1:])
		pkg.Data = ast.NewScope(nil) // required by ast.NewPackage for dot-import
		imports[path] = pkg
	}
	return pkg, nil
}

func parseFiles(fset *token.FileSet, abspath string, names []string) (map[string]*ast.File, error) {
	files := make(map[string]*ast.File)
	for _, f := range names {
		absname := filepath.Join(abspath, f)
		file, err := parser.ParseFile(fset, absname, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[absname] = file
	}
	return files, nil
}

func ImportPackage(ctx *mux.Context, p string) (*Package, error) {
	bctx := buildContext()
	b, err := bctx.Import(p, "", 0)
	if err != nil {
		if noBuildable(err) {
			return nil, err
		}
		b, err = bctx.ImportDir(p, 0)
		if err != nil {
			return nil, err
		}
	}
	fset := token.NewFileSet()
	names := b.GoFiles
	names = append(names, b.CgoFiles...)
	files, err := parseFiles(fset, b.Dir, names)
	if err != nil {
		return nil, err
	}
	// NewPackage will always return errors because it won't
	// resolve builtin types.
	a, _ := ast.NewPackage(fset, files, pkgImporter, nil)
	flags := doc.AllMethods
	if p == "builtin" {
		flags |= doc.AllDecls
	}
	pkg := &Package{
		ctx:  ctx,
		fset: fset,
		bpkg: b,
		apkg: a,
		dpkg: doc.New(a, b.ImportPath, flags),
	}
	sub, err := ImportPackages(ctx, b.Dir)
	if err != nil {
		return nil, err
	}
	pkg.Packages = sub
	return pkg, nil
}

func ImportPackages(ctx *mux.Context, dir string) ([]*Package, error) {
	var pkgs []*Package
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		n := v.Name()
		if n == "test_data" || n[0] == '.' || n[0] == '_' {
			continue
		}
		abs := filepath.Join(dir, n)
		// Follow symlinks
		if st, err := os.Stat(abs); err == nil && st.IsDir() {
			pkg, err := ImportPackage(ctx, abs)
			if err != nil {
				if noBuildable(err) {
					sub, err := ImportPackages(ctx, abs)
					if err != nil {
						return nil, err
					}
					pkgs = append(pkgs, sub...)
					continue
				}
				return nil, err
			}
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}

type packageGroup struct {
	Title    string
	Packages []*Package
}

func DocListHandler(ctx *mux.Context) {
	pkgs, err := ImportPackages(ctx, filepath.Join(srcDir, "gnd.la"))
	if err != nil {
		panic(err)
	}
	groups := []packageGroup{
		{Title: "Gondola Packages", Packages: pkgs},
	}
	infos, _ := ioutil.ReadDir(srcDir)
	var opkgs []*Package
	for _, v := range infos {
		if v.IsDir() && v.Name() != "gnd.la" {
			p, _ := ImportPackages(ctx, filepath.Join(srcDir, v.Name()))
			opkgs = append(opkgs, p...)
		}
	}
	if len(opkgs) > 0 {
		groups = append(groups, packageGroup{
			Title:    "Other useful Rainy Cape's packages",
			Packages: opkgs,
		})
	}
	title := "Package Index"
	data := map[string]interface{}{
		"HeaderTitle": title,
		"Subtitle":    title,
		"Section":     "docs",
		"Groups":      groups,
	}
	ctx.MustExecute("pkgs.html", data)
}

func DocHandler(ctx *mux.Context) {
	path := ctx.IndexValue(0)
	pkg, err := ImportPackage(ctx, path)
	if err != nil {
		if noBuildable(err) {
			sub, err := ImportPackages(ctx, filepath.Join(srcDir, path))
			if err != nil {
				panic(err)
			}
			pkg = &Package{name: path, Packages: sub}
		} else {
			panic(err)
		}
	}
	title := "Package " + pkg.ImportPath()
	breadcrumbs := []*Breadcrumb{
		{Title: "Index", Href: ctx.MustReverse("doc-list")},
	}
	for ii := 0; ii < len(path); {
		var end int
		slash := strings.IndexByte(path[ii:], '/')
		if slash < 0 {
			end = len(path)
		} else {
			end = ii + slash
		}
		breadcrumbs = append(breadcrumbs, &Breadcrumb{
			Title: path[ii:end],
			Href:  ctx.MustReverse("doc", path[:end]),
		})
		ii = end + 1
	}
	data := map[string]interface{}{
		"HeaderTitle": "Package " + pkg.Name(),
		"Subtitle":    title,
		"Section":     "docs",
		"Breadcrumbs": breadcrumbs,
		"Package":     pkg,
	}
	ctx.MustExecute("package.html", data)
}
