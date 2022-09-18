package annotation

import (
	"go/ast"
	"strings"
)

type Annotation struct {
	Key   string
	Value string
}

type annotations struct {
	Node ast.Node
	Ans  []Annotation
}

func newAnnotations(n ast.Node, cg *ast.CommentGroup) annotations {
	if cg == nil || len(cg.List) == 0 {
		return annotations{Node: n}
	}
	ans := make([]Annotation, 0, len(cg.List))
	for _, c := range cg.List {
		text, ok := extractContent(c)
		if !ok {
			continue
		}
		key, value, ok := parseKeyValue(text)
		if !ok {
			continue
		}
		ans = append(ans, Annotation{Key: key, Value: value})
	}
	return annotations{Node: n, Ans: ans}
}

func extractContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	}
	if strings.HasPrefix(text, "/* ") {
		return text[3 : len(text)-2], true
	}
	return "", false
}

func parseKeyValue(text string) (string, string, bool) {
	if !strings.HasPrefix(text, "@") {
		return "", "", false
	}
	segments := strings.SplitN(text, " ", 2)
	if len(segments) != 2 {
		return "", "", false
	}
	key := segments[0][1:]
	value := segments[1]
	return key, value, true
}
