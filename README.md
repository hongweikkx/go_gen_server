[![Go Report Card](https://goreportcard.com/badge/github.com/hongweikkx/go_pattern)](https://goreportcard.com/report/github.com/hongweikkx/go_pattern)
# go_pattern
A Go library that implements groutine work pattern

+ [x] publish/subscribe server
+ [x] parallel (exec task in parallel)
+ [x] Work scheduler (do n same task use m goroutine)
+ [x] actor (actor model)

## Usage
parallel model example(others can refer to the test file):
```
import (
	"fmt"
	pattern "github.com/hongweikkx/go_pattern"
)

func main() {
	p1 := pattern.NewParaller()
	addRet := 0
	var convertA string
	var convertB bool
	p1.Add(pattern.NewHandlerFunc(add, 1, 2).SetRets(&addRet))
	p1.Add(pattern.NewHandlerFunc(convert, "hello", false).SetRets(&convertA, &convertB))
	err := p1.Run()
	if err != nil {
		fmt.Printf("p1 test err:%s\n", err.Error())
		return
	}
	fmt.Printf("addRet:%+v, convertA:%+v, convertB:%+v\n", addRet, convertA, convertB)
}

func add(x, y int) int {
	return x + y
}

func convert(a interface{}, b interface{}) (string, bool) {
	return a.(string), b.(bool)
}
```


## Test
`go test ./...`

