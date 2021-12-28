package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	fName := flag.String("file", "./example/main.go", "./example/main.go")
	flag.Parse()
	fmt.Println(*fName)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *fName, nil, 0)
	if err != nil {
		log.Fatalln(err)
		return
	}

	addWithContext(f)
	format.Node(os.Stdout, fset, f)
}

func addWithContext(f *ast.File) {
	for _, v := range f.Decls {
		switch decl := v.(type) {
		case *ast.GenDecl:
			for _, s := range decl.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}

				switch sst := ts.Type.(type) {
				case *ast.InterfaceType:
					addContextToInterface(sst)
				}
			}
		case *ast.FuncDecl:
			if params := contextFieldList(decl.Type.Params); params != nil {
				declWithContext := &ast.FuncDecl{
					Doc:  nil,
					Recv: decl.Recv,
					Name: &ast.Ident{Name: decl.Name.Name + "WithContext"},
					Type: &ast.FuncType{
						Func:    decl.Type.Func,
						Params:  contextFieldList(decl.Type.Params),
						Results: nil,
					},
					Body: decl.Body,
				}

				f.Decls = append(f.Decls, declWithContext)
			}
		}
	}
}

func addContextToInterface(sst *ast.InterfaceType) {
	q := make([]*ast.Field, 0, 0)
	for _, m := range sst.Methods.List {
		name := m.Names[0].String()
		z := &ast.Field{
			Doc:   m.Doc,
			Names: []*ast.Ident{},
			Type: &ast.FuncType{
				Func: m.Type.(*ast.FuncType).Func,
				Params: &ast.FieldList{
					List: m.Type.(*ast.FuncType).Params.List,
				},
				Results: m.Type.(*ast.FuncType).Results,
			},
			Tag:     m.Tag,
			Comment: m.Comment,
		}

		w, ok := z.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		if len(w.Params.List) == 0 || !isFirstContextContext(w.Params.List[0]) {
			w.Params = contextFieldList(w.Params)
			z.Names = append(z.Names, &ast.Ident{
				NamePos: 0,
				Name:    name + "WithContext",
				Obj:     nil,
			})
			q = append(q, z)
		}
	}
	sst.Methods.List = append(sst.Methods.List, q...)
}

func contextFieldList(list *ast.FieldList) *ast.FieldList {
	if len(list.List) == 0 || !isFirstContextContext(list.List[0]) {
		l := &ast.FieldList{
			List: []*ast.Field{
				{
					Doc: nil,
					Names: []*ast.Ident{
						{
							Name: "ctx",
							Obj: &ast.Object{
								Kind: ast.Var,
								Name: "ctx",
							},
						},
					},
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "context",
						},
						Sel: &ast.Ident{
							Name: "Context",
						},
					},
				},
			},
		}

		if len(list.List) != 0 {
			for _, v := range list.List {
				l.List = append(l.List, v)
			}
		}

		return l
	}

	return nil
}

func isFirstContextContext(spec *ast.Field) bool {
	t := spec.Type
	switch u := t.(type) {
	case *ast.SelectorExpr:
		if u.X.(*ast.Ident).Name == "context" && u.Sel.Name == "Context" {
			return true
		}
	}
	return false
}
