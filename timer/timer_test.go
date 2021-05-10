package timer

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	const N = 1000000
	const d = time.Second * 1

	fmt.Printf("ready add %d timers\n", N)

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		SetTimeoutFunc(d, func(id ID) {
			wg.Done()
		})
	}
	wg.Wait()
}
