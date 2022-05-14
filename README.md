# Go ThreadLocal

Go version of ThreadLocal is like [Java ThreadLocal](https://docs.oracle.com/javase/8/docs/api/java/lang/ThreadLocal.html)

## Usage

```go
import (
	"fmt"
	"sync"

	"github.com/cyub/threadlocal"
)

var wg sync.WaitGroup

func main() {
	tl1 := threadlocal.New()
	tl1.Set("hello, world")

	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			fmt.Printf("id: %d, tid: %d, val: %v\n", idx, threadlocal.ThreadId(), tl1.Get())
		}(i)
	}
	wg.Wait()
	fmt.Printf("id: %d, tid: %d, val: %v\n", n, threadlocal.ThreadId(), tl1.Get())
}
```