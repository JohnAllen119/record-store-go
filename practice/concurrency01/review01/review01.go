package main
import(
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	wg.Add(3)
	ch:=make(chan int)
	go func(){
		defer wg.Done()
		ch <- 10
	}()
	go func(){
		defer wg.Done()
		ch <- 20
	}()
	go func(){
		defer wg.Done()
		ch <- 30
	}()
	go func(){
		wg.Wait()
		close(ch)
	}()
	for x := range ch{
		fmt.Println("x:",x)
	}
	
}