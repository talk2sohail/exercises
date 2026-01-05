package pipeline

import (
	"fmt"
	"sync"
)

// stage: producer
func Gen(done chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// stage: modifier
func Square(done chan struct{}, ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range ch {
			select {
			case out <- n * n:
			case <-done:
				return
			}
			
		}
	}()
	return out
}


// stage: merge (fan-in pattern)
func Merge(done chan struct{}, chs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(ch <-chan int) {
		defer wg.Done()
		for n := range ch {
			select {
			case out <- n:
			case <-done:
				fmt.Println("exiting from merge..")
				return
			}
		}
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
