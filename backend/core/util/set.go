package util

import (
	"go/types"
	"sync"
)

type Set[T comparable] struct {
	elements map[T]*types.Nil
	mut      sync.Mutex
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
		mut:      sync.Mutex{},
	}
}

func (s *Set[T]) Add(value T) bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	if contains := s.Contains(value); contains {
		return false
	}

	s.elements[value] = nil
	return true
}

func (s *Set[T]) Remove(value T) bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	if notContains := s.NotContains(value); notContains {
		return false
	}

	delete(s.elements, value)
	return true
}

func (s *Set[T]) Contains(value T) bool {
	_, exists := s.elements[value]
	return exists
}

func (s *Set[T]) NotContains(value T) bool {
	_, exists := s.elements[value]
	return !exists
}

func (s *Set[T]) Len() int {
	return len(s.elements)
}
