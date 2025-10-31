package annotation

import "github.com/lindb/lindb/sql/tree"

var annotations = make(map[string]CraeteAnnotation)

type Annotation interface{}

type CraeteAnnotation func(ananotion *tree.Annotation) Annotation

func Register(name string, create CraeteAnnotation) {
	if _, ok := annotations[name]; ok {
		panic("annotation exsit")
	}
	annotations[name] = create
}
