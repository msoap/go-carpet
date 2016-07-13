package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Func - one go function in source with return information
type Func struct {
	Name       string
	Begin, End int
}

// getGolangFuncs - parse golang source file and get all functions
//   funcs, err := getGolangFuncs(goFileContentInBytes)
func getGolangFuncs(fileContent []byte) (result []Func, err error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", fileContent, 0)
	if err != nil {
		return result, err
	}

	ast.Inspect(astFile, func(nodeRaw ast.Node) bool {
		switch node := nodeRaw.(type) {
		case *ast.FuncDecl:
			result = append(result, Func{
				Name:  node.Name.String(),
				Begin: int(node.Pos()),
				End:   int(node.End()),
			})
		}

		return true
	})

	return result, nil
}
