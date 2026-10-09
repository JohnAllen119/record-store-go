package main
import(
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0
	wg.Add(1000)
	for i := 0;i<1000;i++{
		go func(){
			mu.Lock()
			count++
			mu.Unlock()
			wg.Done()
		}()
		
	}
	wg.Wait()
	fmt.Println("Count=",count)
}