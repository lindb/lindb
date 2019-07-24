package collections

// Stack represents a last-in-first-out(LIFO) stack of interface, using slice saving data.
// NOTICE: not safe for goroutine concurrent
type Stack struct {
	elements []interface{}
	len      int
}

// NewStack creates an empty stack
func NewStack() *Stack {
	return &Stack{}
}

// Push pushes an item onto the top of this stack
func (s *Stack) Push(item interface{}) {
	s.elements = append(s.elements, item)
	s.len++
}

// Pop removes the item at the top of this stack, if stack is empty return nil
func (s *Stack) Pop() interface{} {
	if s.len == 0 {
		return nil
	}
	item := s.Peek()

	s.elements = s.elements[0 : s.len-1]
	s.len--
	return item
}

// Peek looks at the item at the top of this stack without removing it from the stack
func (s *Stack) Peek() interface{} {
	if s.len == 0 {
		return nil
	}
	return s.elements[s.len-1]
}

// Empty tests if this stack is empty
func (s *Stack) Empty() bool {
	return s.len == 0
}
