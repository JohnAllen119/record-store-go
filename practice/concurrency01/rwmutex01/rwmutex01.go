package main
import(
	"fmt"
	"sync"
	"time"
)

func main(){
	var mu sync.RWMutex
	var wg sync.WaitGroup
	x := 0
	wg.Add(3)
	go func(){
		defer wg.Done()
		mu.RLock()
		fmt.Println("开始读")
time.Sleep(2 * time.Second)
fmt.Println(x)
		mu.RUnlock()
	}()
	go func(){
		defer wg.Done()
		mu.RLock()
		fmt.Println("开始读")
time.Sleep(2 * time.Second)
fmt.Println(x)
		mu.RUnlock()
	}()
	go func(){
		defer wg.Done()
		mu.Lock()
		x++
		mu.Unlock()
	}()
	wg.Wait()
}