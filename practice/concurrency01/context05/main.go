package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan string)
	go func() {
		defer wg.Done()
		ch <- "started"
		<-ctx.Done()
		ch <- "query canceled"
	}()
	s := <-ch
	fmt.Println(s)
	cancel()
	q := <-ch
	fmt.Println(q)
	wg.Wait()
}
