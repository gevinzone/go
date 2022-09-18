package annotation

import (
	"go/ast"
)

type SingleFileEntryVisitor struct {
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) (w ast.Visitor) {
	file, ok := node.(*ast.File)
	if !ok {
		return s
	}
	s.file = &fileVisitor{ans: newAnnotations(file, file.Doc)}
	return s.file
}

func (s *SingleFileEntryVisitor) Get() File {
	if s.file == nil {
		return File{}
	}
	return s.file.Get()
}

type fileVisitor struct {
	ans     annotations
	types   []*typeVisitor
	visited bool
}

func (f *fileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	typ, ok := node.(*ast.TypeSpec)
	if !ok {
		return f
	}
	res := &typeVisitor{
		ans:    newAnnotations(typ, typ.Doc),
		fields: make([]Field, 0, 0),
	}
	f.types = append(f.types, res)
	return res
}

func (f *fileVisitor) Get() File {
	types := make([]Type, 0, len(f.types))
	for _, t := range f.types {
		types = append(types, t.Get())
	}
	return File{
		annotations: f.ans,
		Types:       types,
	}
}

type typeVisitor struct {
	ans     annotations
	fields  []Field
	visited bool
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fd, ok := node.(*ast.Field)
	if !ok {
		return t
	}
	t.fields = append(t.fields, Field{
		annotations: newAnnotations(fd, fd.Doc),
	})
	return nil
}

func (t *typeVisitor) Get() Type {
	return Type{
		annotations: t.ans,
		Fields:      t.fields,
	}
}

type File struct {
	annotations
	Types []Type
}

type Type struct {
	annotations
	Fields []Field
}

type Field struct {
	annotations
}
