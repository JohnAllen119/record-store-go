package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		fmt.Println(ctx.Err())
	}()
	wg.Wait()
}
