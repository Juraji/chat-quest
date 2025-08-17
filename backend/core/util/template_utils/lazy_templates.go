package template_utils

import (
	"fmt"
	"text/template"
)

type LazyTemplateVar[T any] struct {
	loadFunc func() (T, error)
	value    *T
	loaded   bool
}

func NewLazyTemplateVar[T any](loadFunc func() (T, error)) *LazyTemplateVar[T] {
	return &LazyTemplateVar[T]{
		loadFunc: loadFunc,
	}
}

func (p *LazyTemplateVar[T]) Get() (*T, error) {
	if !p.loaded {
		value, err := p.loadFunc()
		if err != nil {
			return nil, err
		}

		p.value = &value
		p.loaded = true
	}

	return p.value, nil
}

func (p *LazyTemplateVar[T]) String() string {
	get, err := p.Get()
	if err != nil {
		panic(err)
	}
	if get == nil {
		return ""
	}

	return template.HTMLEscapeString(fmt.Sprint(get))
}

type LazyTemplateSlice[T any] struct {
	loadFunc func() ([]T, error)
	value    []T
	loaded   bool
}

func NewLazyTemplateSlice[T any](loadFunc func() ([]T, error)) *LazyTemplateSlice[T] {
	return &LazyTemplateSlice[T]{
		loadFunc: loadFunc,
	}
}

func (p *LazyTemplateSlice[T]) Get() ([]T, error) {
	if !p.loaded {
		value, err := p.loadFunc()
		if err != nil {
			return nil, err
		}

		p.value = value
		p.loaded = true
	}

	return p.value, nil
}

func (p *LazyTemplateSlice[T]) String() string {
	get, err := p.Get()
	if err != nil {
		panic(err)
	}
	if len(get) == 0 {
		return ""
	}
	return fmt.Sprint(len(get))
}

func LazyTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"getLazy": func(l *LazyTemplateVar[any]) (any, error) {
			return l.Get()
		},
		"getLazySlice": func(l *LazyTemplateSlice[any]) ([]any, error) {
			return l.Get()
		},
	}
}
