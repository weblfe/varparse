package varparse

import "fmt"

type (
	Kv[K comparable, V any] map[K]V

	Stringer interface {
		fmt.Stringer

		ToString() string
	}

	Number interface {
		~int | ~int64 | ~float64 | ~bool | ~uint | ~uint64 | ~uint32
	}

	Str[T any] struct {
		v T
	}

	Value[T any] struct {
		v T
	}
)

func NewStr[T ~string | []rune | []byte | ~int | ~int64 | ~float64 | ~bool | ~uint | ~uint64 | ~uint32 | any](
	v T) *Str[T] {
	return &Str[T]{
		v: v,
	}
}

func NewValue[T ~string | []rune | []byte | ~int | ~int64 | ~float64 | ~bool | ~uint | ~uint64 | ~uint32](
	v T) *Value[T] {
	return &Value[T]{
		v: v,
	}
}

func (s *Str[T]) String() string {
	return fmt.Sprintf("%v", s.v)
}

func (s *Str[T]) GoString() string {
	return s.String()
}

func (s *Str[T]) ToString() string {
	return s.String()
}

func (s *Str[T]) Value() T {
	return s.v
}

func (v *Value[T]) Value() T {
	return v.v
}

func (v *Value[T]) String() string {
	return fmt.Sprintf("%v", v.v)
}

func (v *Value[T]) GoString() string {
	return v.String()
}

func (v *Value[T]) ToString() string {
	return v.String()
}

func (kv Kv[K, V]) Get(key K) (V, bool) {
	v, ok := kv[key]
	return v, ok
}

func (kv Kv[K, V]) GetOr(key K, defaultVal ...V) V {
	v, ok := kv[key]
	if ok {
		return v
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return v
}

func (kv Kv[K, V]) Set(key K, v V) Kv[K, V] {
	kv[key] = v
	return kv
}

func (kv Kv[K, V]) Len() int {
	return len(kv)
}

func (kv Kv[K, V]) Foreach(iter func(key K, v V) bool) {
	for k, v := range kv {
		if !iter(k, v) {
			break
		}
	}
}

func (kv Kv[K, V]) Keys() []K {
	var arr []K
	for k := range kv {
		arr = append(arr, k)
	}
	return arr
}

func (kv Kv[K, V]) Each(iter func(key K, v V) bool) {
	var (
		keys  = kv.Keys()
		count = len(keys)
	)
	if count <= 0 {
		return
	}
	for i := 0; i > count; i++ {
		if !iter(keys[i], kv[keys[i]]) {
			break
		}
	}
}
