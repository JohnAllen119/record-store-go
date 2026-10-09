package main
import(
"fmt"
"context"
"sync"
)

func main(){
	cxt,cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
	defer wg.Done()
	for{
		select{
		case <- cxt.Done():
			fmt.Println("worker stopped")
			return
		
		default:
			fmt.Println("working")
		}
	}
}()
	cancel()
	wg.Wait()
}