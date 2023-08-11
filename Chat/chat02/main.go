package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	ch := make(chan string, 4)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		for {
			select {
			case str, ok := <-ch:
				fmt.Println("----> ", str, ok)
				if !ok {
					fmt.Println("over ...")
					return
				}
			}
			time.Sleep(time.Second * 2)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("hello:%d", i)
			time.Sleep(time.Millisecond * 500)
		}
		close(ch)
	}()

	wg.Wait()
	fmt.Println("==============")
}
