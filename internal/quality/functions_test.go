package quality

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestFunctionsStayWithinThirtyLines(t *testing.T) {
	err := filepath.WalkDir("../..", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".." {
			return filepath.SkipDir
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".go") {
			checkFunctions(t, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func checkFunctions(t *testing.T, path string) {
	t.Helper()
	positions := token.NewFileSet()
	file, err := parser.ParseFile(positions, path, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(file, func(node ast.Node) bool {
		var start, end token.Pos
		name := "anonymous function"
		switch function := node.(type) {
		case *ast.FuncDecl:
			start, end, name = function.Pos(), function.End(), function.Name.Name
		case *ast.FuncLit:
			start, end = function.Pos(), function.End()
		default:
			return true
		}
		lines := positions.Position(end).Line - positions.Position(start).Line + 1
		if lines > 30 {
			t.Errorf("%s: %s is %d lines; hard limit is 30", path, name, lines)
		}
		return true
	})
}
