package varparse_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/weblfe/varparse"
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	var p = varparse.NewParser[string, fmt.Stringer]()
	p.Assign("test", varparse.NewValue("test你好"))
	p.Assign("number", varparse.NewValue(1))
	p.Assign("bool", varparse.NewValue(true))
	extractor := func(s string) map[string]string {
		if !strings.Contains(s, "${number}") && !strings.Contains(s, "${bool}") {
			return nil
		}
		return map[string]string{
			"number": "${number}",
			"bool":   "${bool}",
		}
	}
	result := p.Parse("${number}/test/${bool}", extractor)
	var expected = "1/test/true"
	t.Logf("%v", result)
	t.Logf("%v", p.Get("test"))
	var as = assert.New(t)
	as.Equal(expected, result, "parse failed")
}

func TestNewExtractor(t *testing.T) {
	type caseItem struct {
		value  string
		start  string
		end    string
		expect map[string]string
	}
	var as = assert.New(t)
	var cases = []caseItem{
		{
			value: `${number}/test/${bool}`,
			start: "${",
			end:   "}",
			expect: map[string]string{
				"number": "${number}",
				"bool":   "${bool}",
			},
		},
		{
			value: `${number}/test/${bool}/${/ssx}`,
			start: "${",
			end:   "}",
			expect: map[string]string{
				"number": "${number}",
				"bool":   "${bool}",
			},
		},
		{
			value: `$number/test/$bool/$/ssx`,
			start: "$",
			end:   "/",
			expect: map[string]string{
				"number": "$number/",
				"bool":   "$bool/",
			},
		},
		{
			value: `<number>/test/<bool>/</ssx>`,
			start: "<",
			end:   ">",
			expect: map[string]string{
				"number": "<number>",
				"bool":   "<bool>",
			},
		},
		{
			value: `:number/test/:b_id/:/ssx:`,
			start: ":",
			end:   "/",
			expect: map[string]string{
				"number": ":number/",
				"b_id":   ":b_id/",
			},
		},
		{
			value: `[:number]/test/[:b_id]/[:/ssx]`,
			start: "[:",
			end:   `]`,
			expect: map[string]string{
				"number": "[:number]",
				"b_id":   "[:b_id]",
			},
		},
	}
	for _, v := range cases {
		t.Run("-"+v.value, func(t *testing.T) {
			executor := varparse.NewExtractor(v.start, v.end)
			err := executor.Compile()
			if err != nil {
				t.Error(err)
				return
			}
			values := executor.Extract(v.value)
			for k, v := range values {
				t.Logf("%s=>%s", k, v)
			}
			as.Equal(v.expect, values, "extract fail")
		})
	}
}

func TestParseImpl_Parse(t *testing.T) {
	type caseItem struct {
		value  string
		start  string
		end    string
		vars   map[string]any
		expect string
	}
	var as = assert.New(t)
	var cases = []caseItem{
		{
			value: `${number}/test/${bool}`,
			start: "${",
			end:   "}",
			vars: map[string]any{
				"number": 123,
				"bool":   true,
			},
			expect: "123/test/true",
		},
		{
			value: `${number}/test/${bool}/${/ssx}`,
			start: "${",
			end:   "}",
			vars: map[string]any{
				"number": "0001",
				"bool":   false,
			},
			expect: "0001/test/false/${/ssx}",
		},
		{
			value: `$number/test/$bool/$/ssx`,
			start: "$",
			end:   "/",
			vars: map[string]any{
				"number": "12001",
				"bool":   "0",
			},
			expect: "12001test/0$/ssx",
		},
		{
			value: `<number>/test/<bool>/</ssx>`,
			start: "<",
			end:   ">",
			vars: map[string]any{
				"number": "nu_001",
				"bool":   false,
			},
			expect: "nu_001/test/false/</ssx>",
		},
		{
			value: `:number/test/:b_id/:/ssx:`,
			start: ":",
			end:   "/",
			vars: map[string]any{
				"number": "${id}",
				"b_id":   "001",
				"id":     123,
			},
			expect: "${id}test/001:/ssx:",
		},
		{
			value: `${number}/test/${b_id}/:/ssx:`,
			start: "${",
			end:   "}",
			vars: map[string]any{
				"number": "${id}",
				"b_id":   "001",
				"id":     123,
			},
			expect: "123/test/001/:/ssx:",
		},
		{
			value: `[:number]/test/[:b_id]/[:/ssx]`,
			start: "[:",
			end:   `]`,
			vars: map[string]any{
				"number": "[000]",
				"b_id":   "111",
			},
			expect: "[000]/test/111/[:/ssx]",
		},
	}
	var p = varparse.NewParser[string, fmt.Stringer]()

	for _, v := range cases {
		for k, vx := range v.vars {
			p.Assign(k, varparse.NewStr(vx))
		}
		t.Run("-"+v.value, func(t *testing.T) {
			executor := varparse.NewExtractor(v.start, v.end)
			err := executor.Compile()
			if err != nil {
				t.Error(err)
				return
			}
			content := p.Parse(v.value, executor.Extract)
			t.Logf("content=%v", content)
			as.Equal(v.expect, content, "extract parse fail")
		})
	}
}
