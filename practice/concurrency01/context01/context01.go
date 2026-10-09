package main
import(
	"context"
	"fmt"
	"sync"
)
func main(){
	ctx,cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		defer wg.Done()
		<-ctx.Done()
		fmt.Println("worker stopped")
	}()
	cancel()
	wg.Wait()
}