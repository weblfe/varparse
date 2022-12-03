package varparse

import (
	"fmt"
	"strings"
)

type VarParse[K comparable, V any] interface {
	Assign(K, V) VarParse[K, V]
	Get(K) V
	Parse(content string, extract ExtractHandler[K]) string
}

type ExtractHandler[K comparable] func(string) map[K]string

type parseImpl[K comparable, V any] struct {
	vars Kv[K, V]
}

func (p *parseImpl[K, V]) Assign(k K, v V) VarParse[K, V] {
	p.vars.Set(k, v)
	return p
}

func (p *parseImpl[K, V]) Get(k K) V {
	return p.vars.GetOr(k)
}

func NewParser[K comparable, V any]() VarParse[K, V] {
	var impl = &parseImpl[K, V]{
		vars: make(Kv[K, V]),
	}
	return impl
}

func (p *parseImpl[K, V]) Parse(content string, extract ExtractHandler[K]) string {
	values := extract(content)
	for k, vk := range values {
		vx, ok := p.vars[k]
		if !ok {
			continue
		}
		content = strings.ReplaceAll(content, vk, fmt.Sprintf("%v", vx))
	}
	return content
}
