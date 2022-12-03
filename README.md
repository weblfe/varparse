# a var parse go lib

intro
   
   A parse dynamic variable templates lib

## install
```shell
go get github.com/weblfe/varparse
```

## example

```go
package main

import (
	"fmt"
	"github.com/weblfe/varparse"
)

func main() {
		var p = varparse.NewParser[string, any]()
		p.Assign("en", varparse.NewValue("hello,Go"))
		p.Assign("cn", varparse.NewValue("你好Go"))
		p.Assign("bool", varparse.NewValue(true))
		executor := varparse.NewExtractor("${","}")
		// default:=varparse.ExtractorOf()
		err := executor.Compile()
		if err!=nil {
				panic(err)
        }
		content := p.Parse(`${en},${cn} this a lib for var parse! Yes! ${bool}`, executor.Extract)		
		fmt.Println(content)
}
```
