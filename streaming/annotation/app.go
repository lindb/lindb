package annotation

import "github.com/lindb/lindb/sql/tree"

func init() {
	Register("app", func(ananotion *tree.Annotation) Annotation {
		return &AppAnnotation{}
	})
}

type AppAnnotation struct {
	Name string
}
