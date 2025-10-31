package tree

type AnnotationElement interface{}

type Annotation struct {
	Name     *Identifier
	Elements []AnnotationElement
}
