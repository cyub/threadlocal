package threadlocal

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestGid(t *testing.T) {
	var wg sync.WaitGroup
	n := 100
	wg.Add(100)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			i := Gid()
			j := extractGidFromStack()
			if i != j {
				t.Fatal(fmt.Sprintf("%d != %d", i, j))
			}
		}()
	}
	wg.Wait()
}

func extractGidFromStack() int {
	var buf [128]byte
	n := runtime.Stack(buf[:], false)
	stack := strings.Split(string(buf[:n]), " ")
	id, _ := strconv.Atoi(stack[1])
	return id
}
