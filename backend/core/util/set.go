package util

import (
	"go/types"
)

type Set[T comparable] struct {
	elements map[T]*types.Nil
}

func NewSetFrom[T any, K comparable](slice []T, transform func(t T) K) *Set[K] {
	set := NewSet[K](len(slice))

	for _, value := range slice {
		set.Add(transform(value))
	}

	return set
}

func NewSet[T comparable](initialSize int) *Set[T] {
	return &Set[T]{
		elements: make(map[T]*types.Nil, initialSize),
	}
}

func (s *Set[T]) Add(value T) {
	s.elements[value] = nil
}

func (s *Set[T]) Remove(value T) {
	delete(s.elements, value)
}

func (s *Set[T]) Contains(value T) bool {
	_, ok := s.elements[value]
	return ok
}

func (s *Set[T]) NotContains(value T) bool {
	_, ok := s.elements[value]
	return !ok
}

func (s *Set[T]) Len() int {
	return len(s.elements)
}
