package main
import(
	"fmt"
	"time"
)

func main(){
	ch2 := make(chan int)
	ch1 := make(chan int)
	go func(){
		time.Sleep(1*time.Second)
		ch1 <- 10
	}()
	go func(){
		time.Sleep(2*time.Second)
		ch2 <- 20
	}()
	for i := 0; i < 2; i++ {
	select{
	case x := <- ch1:
		fmt.Println("收到:",x)
	case y := <- ch2:
		fmt.Println("收到:",y)
	case <-time.After(3*time.Second):
		fmt.Println("超时")
	}
}
}