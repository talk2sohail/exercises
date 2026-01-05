package main

import (
	"fmt"
	"os"
	"pipeline/filewalk"
	"sort"
)

// main: for demonstrating the fan-in/fan-out pattern
// func main() {
// 	done := make(chan struct{})
// 	ch := pipeline.Gen(done, 2, 4, 3, 7)
// 	r1 := pipeline.Square(done, ch)
// 	r2 := pipeline.Square(done, ch)
// 	i := 0
// 	for v := range pipeline.Merge(done, r1, r2) {
// 		if i == 1 {
// 			close(done) // send the broadcast to upstreams, not interested with remaining values
// 		}
// 		fmt.Println(v)
// 		i++
// 	}

// }

func main() {

	m, err := filewalk.Parallel_MD5All_Bounded(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}

}
