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

// Push pushes an element onto the top of this stack
func (s *Stack) Push(element interface{}) {
	s.elements = append(s.elements, element)
	s.len++
}

// Pop removes the element at the top of this stack, if stack is empty return nil
func (s *Stack) Pop() interface{} {
	if s.len == 0 {
		return nil
	}
	element := s.Peek()

	s.elements = s.elements[0 : s.len-1]
	s.len--
	return element
}

// Peek looks at the element at the top of this stack without removing it from the stack
func (s *Stack) Peek() interface{} {
	if s.len == 0 {
		return nil
	}
	return s.elements[s.len-1]
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return s.len
}

// Empty tests if this stack is empty
func (s *Stack) Empty() bool {
	return s.len == 0
}
