package main

import(
	"fmt"
	"sync"
)
func main(){
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(10)
	for i:=0 ; i<10 ; i++{

		go func(){
		defer wg.Done()
		once.Do(func(){
		fmt.Println("初始化数据连接")
	})
	}()
	
	
	
}
	wg.Wait()
	return
}