package processor

type Processor interface {
	Run(output chan<- any)
	Children() []Processor

	GetInbounds() []chan any
	String() string
}
