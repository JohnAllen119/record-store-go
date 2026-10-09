package main
import(
	"fmt"
	"sync"
)
func main(){
	scores := make(map[int]int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i:= 0 ; i<100;i++{
		wg.Add(1)
		go func(id int){
			defer wg.Done()

			mu.Lock()
			scores[id] = id *10
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Println(len(scores))
}