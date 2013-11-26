package main

import (
	"errors"
	"github.com/golang/lint"
	"gnd.la/util/astutil"
	"go/ast"
	"go/doc"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	valueScore = 1
	funcScore  = 2
	typeScore  = 3
	fieldScore = 1
)

var (
	errInvalidPackage = errors.New("invalid package")
)

type UndocumentedKind int

const (
	Const UndocumentedKind = iota + 1
	Var
	Func
	Type
	Method
	Field
	IMethod
)

func (k UndocumentedKind) Score() int {
	switch k {
	case Const, Var:
		return valueScore
	case Func, Method, IMethod:
		return funcScore
	case Type:
		return typeScore
	case Field:
		return fieldScore
	case 0:
		return 0
	}
	panic("unreachable")
}

type Undocumented struct {
	Kind UndocumentedKind
	Name string
	Type string
}

func (u *Undocumented) String() string {
	switch u.Kind {
	case Const:
		return "constant " + u.Name
	case Var:
		return "variable " + u.Name
	case Func:
		return "function " + u.Name
	case Type:
		return "type " + u.Name
	case Method:
		return "method (" + u.Type + ") " + u.Name
	case Field:
		return "field " + u.Name + " on type " + u.Type
	case IMethod:
		return "method " + u.Name + " on interface " + u.Type
	}
	return "invalid Undocumented"
}

func (u *Undocumented) Id() string {
	switch u.Kind {
	case Const:
		return ConstId(u.Name)
	case Var:
		return VarId(u.Name)
	case Func:
		return FuncId(u.Name)
	case Type:
		return TypeId(u.Name)
	case Method:
		return MethodId(u.Type, u.Name)
	case Field:
		return TypeId(u.Type)
	case IMethod:
		return TypeId(u.Type)
	}
	return ""
}

const (
	noDocPenalty = 10
)

type Stats struct {
	Documented int
	ToDocument int
	// Indicates if the package has documentation.
	HasDoc       bool
	Undocumented []*Undocumented
	Problems     []*lint.Problem
}

func (s *Stats) NoDocPenalty() int {
	return noDocPenalty
}

func (s *Stats) Penalty() int {
	p := 0
	if !s.HasDoc {
		p += noDocPenalty
	}
	return p
}

func (s *Stats) coef() float64 {
	return float64(100 - s.Penalty())
}

func (s *Stats) Score() float64 {
	return s.coef() * float64(s.Documented) / float64(s.ToDocument)
}

func (s *Stats) Increase(k UndocumentedKind) float64 {
	return s.coef() * float64(k.Score()) / float64(s.ToDocument)
}

func (s *Stats) valueStats(k UndocumentedKind, values []*doc.Value, total *int, score *int) {
	for _, v := range values {
		if v.Doc != "" {
			// There's a comment just before the declaration.
			// Consider all the values documented
			c := len(v.Decl.Specs) * valueScore
			*score += c
			*total += c
		} else {
			// Check every value declared in this group
			for _, spec := range v.Decl.Specs {
				*total += valueScore
				sp := spec.(*ast.ValueSpec)
				if sp.Doc != nil || sp.Comment != nil {
					*score += valueScore
				} else {
					for _, n := range sp.Names {
						s.Undocumented = append(s.Undocumented, &Undocumented{
							Kind: k,
							Name: astutil.Ident(n),
						})
					}
				}
			}
		}
	}
}

func (s *Stats) funcStats(typ string, fns []*doc.Func, total *int, score *int) {
	for _, v := range fns {
		// Skip Error() and String() methods
		if typ != "" && (v.Name == "String" || v.Name == "Error") {
			continue
		}
		*total += funcScore
		if v.Doc != "" {
			*score += funcScore
		} else {
			und := &Undocumented{
				Kind: Func,
				Name: v.Name,
			}
			if typ != "" {
				und.Type = typ
				und.Kind = Method
			}
			s.Undocumented = append(s.Undocumented, und)
		}
	}
}

func (s *Stats) typeStats(typs []*doc.Type, total *int, score *int) {
	for _, v := range typs {
		*total += typeScore
		if v.Doc != "" {
			*score += typeScore
		} else {
			s.Undocumented = append(s.Undocumented, &Undocumented{
				Kind: Type,
				Name: v.Name,
			})
		}
		// Fields
		var k UndocumentedKind
		ts := v.Decl.Specs[0].(*ast.TypeSpec)
		var fields []*ast.Field
		switch s := ts.Type.(type) {
		case *ast.StructType:
			fields = s.Fields.List
			k = Field
		case *ast.InterfaceType:
			fields = s.Methods.List
			k = IMethod
		}
		fs := k.Score()
		for _, f := range fields {
			*total += fs
			if f.Doc != nil || f.Comment != nil {
				*score += fs
			} else {
				var name string
				if len(f.Names) > 0 {
					name = astutil.Ident(f.Names[0])
				} else {
					// Embedded field
					name = astutil.Ident(f.Type)
					if name[0] == '*' {
						name = name[1:]
					}
					if dot := strings.IndexByte(name, '.'); dot >= 0 {
						name = name[dot+1:]
					}
				}
				s.Undocumented = append(s.Undocumented, &Undocumented{
					Kind: k,
					Name: name,
					Type: v.Name,
				})
			}
		}
		s.valueStats(Const, v.Consts, total, score)
		s.valueStats(Var, v.Vars, total, score)
		s.funcStats("", v.Funcs, total, score)
		s.funcStats(v.Name, v.Methods, total, score)
	}
}

func (s *Stats) lint(p *Package) ([]*lint.Problem, error) {
	var problems []*lint.Problem
	linter := &lint.Linter{}
	for _, v := range p.Filenames() {
		if filepath.Ext(v) == ".go" {
			path := filepath.Join(p.bpkg.Dir, v)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			pro, err := linter.Lint(v, data)
			if err != nil {
				return nil, err
			}
			for _, pr := range pro {
				// We're already reporting these kind of issues with
				// the doc stats.
				if !strings.Contains(pr.Text, "should have comment") && !strings.Contains(pr.Text, "should have a package comment") {
					pr.Position.Filename = v
					// pr will get overwriten in next iteration
					prc := pr
					problems = append(problems, &prc)
				}
			}
		}
	}
	return problems, nil
}

func NewStats(p *Package) (*Stats, error) {
	if p.dpkg == nil {
		return nil, errInvalidPackage
	}
	s := new(Stats)
	pr, err := s.lint(p)
	if err != nil {
		return nil, err
	}
	s.Problems = pr
	total := 0
	score := 0
	s.HasDoc = p.dpkg.Doc != ""
	s.valueStats(Const, p.dpkg.Consts, &total, &score)
	s.valueStats(Var, p.dpkg.Vars, &total, &score)
	s.funcStats("", p.dpkg.Funcs, &total, &score)
	s.typeStats(p.dpkg.Types, &total, &score)
	s.Documented = score
	s.ToDocument = total
	return s, nil
}
