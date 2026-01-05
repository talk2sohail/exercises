package pipeline

import (
	"time"
)

func ExecuteAfter(delay time.Duration, fn func() error) chan error {
	// call this named fn after the given delay
	t := time.NewTimer(delay)
	errc := make(chan error, 1)

	go func() {
		<-t.C
		errc <- fn()
	}()

	return errc
}
