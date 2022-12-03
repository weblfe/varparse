package varparse

import (
	"errors"
	"fmt"
	"regexp"
)

type extractor struct {
	varStart string
	varEnd   string
	executor *regexp.Regexp
}

type Extractor interface {
	Extract(v string) map[string]string
	Compile() (err error)
}

var (
	MissStartFlagError = errors.New(`miss start flag`)
	MissEndFlagError   = errors.New(`miss end flag`)
)

func NewExtractor(start, end string) Extractor {
	return &extractor{varStart: start, varEnd: end}
}

func ExtractorOf() Extractor {
	return NewExtractor("${", "}")
}

func (e *extractor) Compile() (err error) {
	if e.varStart == "" {
		return MissStartFlagError
	}
	if e.varEnd == "" {
		return MissEndFlagError
	}
	var (
		end   = RegexpPrevProcess([]rune(e.varEnd))
		start = RegexpPrevProcess([]rune(e.varStart))
		expr  = fmt.Sprintf(`(%s([\w_-]+)%s)`, string(start), string(end))
	)
	e.executor, err = regexp.Compile(expr)
	return err
}

func (e *extractor) Extract(v string) map[string]string {
	if e.executor == nil {
		if err := e.Compile(); err != nil {
			return map[string]string{}
		}
	}
	var (
		extractValue = make(map[string]string)
		values       = e.executor.FindAllStringSubmatch(v, -1)
	)
	for _, v := range values {
		extractValue[v[len(v)-1]] = v[0]
	}
	return extractValue
}

func RegexpPrevProcess(word []rune) []rune {
	if len(word) <= 0 {
		return word
	}
	var (
		result    []rune
		metaWords = map[rune]struct{}{
			'?': {}, '$': {},
			'{': {}, '}': {},
			'\\': {}, '.': {},
			'*': {}, '+': {},
			'[': {}, ']': {},
			'|': {},
		}
	)
	for _, v := range word {
		if _, ok := metaWords[v]; !ok {
			result = append(result, v)
		} else {
			result = append(result, '\\', v)
		}
	}
	return result
}
