package main
import(
	"fmt"
	"sync"
)
func worker(ch chan int,a int,wg *sync.WaitGroup){
	ch <- a*2
	wg.Done()
}
func main(){
	var wg sync.WaitGroup
	ch:=make(chan int)
	wg.Add(3)
	go worker(ch,10,&wg)
	go worker(ch,20,&wg)
	go worker(ch,30,&wg)
	go func(){
		wg.Wait()
		close(ch)
	}()
	
	for x := range ch{
		fmt.Println(x)
	} 
}