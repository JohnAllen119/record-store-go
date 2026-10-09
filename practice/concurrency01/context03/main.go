package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				fmt.Println("worker stopped")
				return
			case <-ticker.C:
				fmt.Println("working")
			}
		}
	}()
	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()
}
