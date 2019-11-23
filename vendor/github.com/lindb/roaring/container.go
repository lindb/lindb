package roaring

/**********************************************************************/
// customize for LinDB
/**********************************************************************/

// Container represents the wrapper interface of roaring bitmap's container,
// like array container/run container/bitmap container
type Container interface {
	// And computes the intersection between two bitmap's container and returns the new container
	And(other Container) Container
	// Or computes the union between two bitmaps's container and returns the new container
	Or(other Container) Container
	// PeekableIterator creates a new ShortPeekable to iterate over the shorts contained in the container
	PeekableIterator() PeekableShortIterator
	// GetCardinality returns the number of shorts contained in the container
	GetCardinality() int
	// Contains returns true if the short is contained in the container
	Contains(key uint16) bool
	// Rank returns the number of short that are smaller or equal to x (Rank(infinity) would be GetCardinality())
	Rank(key uint16) int
	// Minimum get the smallest value stored in this container, assumes that it is not empty
	Minimum() uint16
	// Maximum get the largest value stored in this container, assumes that it is not empty
	Maximum() uint16
	// ToArray creates a new slice containing all of the shorts stored in the container in sorted order
	ToArray() []uint16

	// getContainer gets the real container
	getContainer() container
}

// containerWrapper represents the container wrapper
type containerWrapper struct {
	container container
}

func (c *containerWrapper) And(other Container) Container {
	result := c.container.and(other.getContainer())
	return &containerWrapper{container: result}
}

func (c *containerWrapper) Or(other Container) Container {
	result := c.container.or(other.getContainer())
	return &containerWrapper{container: result}
}

func (c *containerWrapper) PeekableIterator() PeekableShortIterator {
	return newPeekableShortIterator(c.container.getShortIterator())
}

func (c *containerWrapper) GetCardinality() int {
	return c.container.getCardinality()
}

func (c *containerWrapper) Rank(key uint16) int {
	return c.container.rank(key)
}

func (c *containerWrapper) Contains(key uint16) bool {
	return c.container.contains(key)
}

func (c *containerWrapper) Minimum() uint16 {
	return c.container.minimum()
}

func (c *containerWrapper) Maximum() uint16 {
	return c.container.maximum()
}

func (c *containerWrapper) ToArray() []uint16 {
	result := make([]uint16, c.container.getCardinality())
	it := c.container.getShortIterator()
	idx := 0
	for it.hasNext() {
		result[idx] = it.next()
		idx++
	}
	return result
}

func (c *containerWrapper) getContainer() container {
	return c.container
}
