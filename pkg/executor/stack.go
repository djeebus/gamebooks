package executor

type stack[T any] struct {
	items []T
}

func newStack[T any](items ...T) *stack[T] {
	return &stack[T]{items: items}
}

func (s *stack[T]) peek() T {
	return s.items[len(s.items)-1]
}

func (s *stack[T]) push(t T) {
	s.items = append(s.items, t)
}

func (s *stack[T]) pop() T {
	switch len(s.items) {
	case 0:
		panic("stack is empty")
	case 1:
		item := s.items[0]
		s.items = nil
		return item
	default:
		item := s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
		return item
	}
}

func (s *stack[T]) empty() bool {
	return len(s.items) == 0
}
