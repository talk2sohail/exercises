package filewalk

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type entry struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func Parallel_MD5All_Bounded(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	done := make(chan struct{})
	defer close(done)

	paths, countc, errc := parallel_md5all_bounded(done, root)

	var wg sync.WaitGroup
	numDigestors := 20

	c := make(chan entry)

	for i := 0; i < numDigestors; i++ {
		wg.Go(func() {
			digester(done, paths, c)
		})
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil {
		return nil, err
	}

	fmt.Println("count of files:", <-countc)

	return m, nil
}

func parallel_md5all_bounded(done chan struct{}, root string) (chan string, chan int64, chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	countc := make(chan int64, 1)
	go func() {
		defer close(paths)
		var count int64
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			count++

			select {
			case paths <- path:
			case <-done:
				return errors.New("walk cancelled")
			}
			return nil
		})

		countc <- count
	}()

	return paths, countc, errc
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- entry) {
	for p := range paths {
		data, err := os.ReadFile(p)
		select {
		case c <- entry{p, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

func Parallel_MD5All_Unbounded(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	done := make(chan struct{})
	defer close(done)

	entries, errc := parallel_md5sum_unbounded(done, root)
	// iterate over the channel to get the file path and sum
	for e := range entries {
		if e.err != nil {
			return nil, e.err
		}
		m[e.path] = e.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}

func parallel_md5sum_unbounded(done chan struct{}, root string) (<-chan entry, <-chan error) {
	results := make(chan entry)
	var wg sync.WaitGroup
	errc := make(chan error, 1)
	go func() {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			wg.Go(func() {
				data, err := os.ReadFile(path)
				select {
				case results <- entry{path, md5.Sum(data), err}:
				case <-done:
				}
			})

			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}

		})

		go func() {
			wg.Wait()
			close(results)
		}()

		errc <- err

	}()

	return results, errc

}
