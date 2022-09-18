package annotation

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSingleFileEntryVisitor_Get(t *testing.T) {
	src := `
// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
second line
*/
// @date 2022/04/02
package annotation

type (
	// FuncType is a type
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	FuncType func()
)

type (
	// StructType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	StructType struct {
		// Public is a field
		// @type string
		Public string
	}

	// SecondType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/03
	SecondType struct {
	}
)

type (
	// Interface is a test interface
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/04
	Interface interface {
		// MyFunc is a test func
		// @parameter arg1 int
		// @parameter arg2 int32
		// @return string
		MyFunc(arg1 int, arg2 int32) string

		// second is a test func
		// @return string
		second() string
	}
)
`
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	rootVisitor := &SingleFileEntryVisitor{}
	ast.Walk(rootVisitor, f)
	file := rootVisitor.Get()
	t.Log(file)
}
