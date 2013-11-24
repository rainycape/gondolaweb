package main

import (
	"bytes"
	"gnd.la/mux"
	"gnd.la/util"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	srcDir    = util.RelativePath("src")
	cwd       = util.RelativePath("..")
	srcPrefix = srcDir[len(cwd)+1:]
)

type Package struct {
	fset     *token.FileSet
	name     string
	bpkg     *build.Package
	apkg     *ast.Package
	dpkg     *doc.Package
	Packages []*Package
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
	if p.dpkg != nil {
		f := p.dpkg.Filenames
		files := make([]string, len(f))
		for ii, v := range f {
			files[ii] = filepath.Base(v)
		}
		return files
	}
	return nil
}

func (p *Package) HasDoc() bool {
	return p.dpkg != nil && strings.TrimSpace(doc.Synopsis(p.dpkg.Doc)) != strings.TrimSpace(p.dpkg.Doc)
}

func (p *Package) Doc() *doc.Package {
	return p.dpkg
}

func (p *Package) HTML(text string) template.HTML {
	var buf bytes.Buffer
	doc.ToHTML(&buf, text, nil)
	return template.HTML(buf.String())
}

func (p *Package) HTMLDoc() template.HTML {
	return p.HTML(p.dpkg.Doc)
}

func (p *Package) HTMLDecl(node interface{}) (template.HTML, error) {
	s, err := FormatNode(p.fset, node)
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

func ImportPackage(p string) (*Package, error) {
	ctx := build.Default
	ctx.GOPATH = util.RelativePath(".")
	b, err := ctx.Import(p, "", 0)
	if err != nil {
		b, err = ctx.ImportDir(p, 0)
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
		fset: fset,
		bpkg: b,
		apkg: a,
		dpkg: doc.New(a, b.ImportPath, flags),
	}
	sub, err := ImportPackages(b.Dir)
	if err != nil {
		return nil, err
	}
	pkg.Packages = sub
	return pkg, nil
}

func ImportPackages(dir string) ([]*Package, error) {
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
			pkg, err := ImportPackage(abs)
			if err != nil {
				if strings.Contains(err.Error(), "no buildable") {
					sub, err := ImportPackages(abs)
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
	pkgs, err := ImportPackages(filepath.Join(srcDir, "gnd.la"))
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
			p, _ := ImportPackages(filepath.Join(srcDir, v.Name()))
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
	pkg, err := ImportPackage(path)
	if err != nil {
		if strings.Contains(err.Error(), "no buildable") {
			sub, err := ImportPackages(filepath.Join(srcDir, path))
			if err != nil {
				panic(err)
			}
			pkg = &Package{name: path, Packages: sub}
		} else {
			panic(err)
		}
	}
	title := "Package " + pkg.Name()
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
		"HeaderTitle": title,
		"Subtitle":    title,
		"Section":     "docs",
		"Breadcrumbs": breadcrumbs,
		"Package":     pkg,
	}
	ctx.MustExecute("package.html", data)
}
