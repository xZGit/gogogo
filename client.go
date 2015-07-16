package main


import (
	"./godis"
    "errors"
	"sync"
	"log"

)


var waitgroup sync.WaitGroup

func Afunction(shownum int) {
	h := func(v godis.ProtoType) (interface{}, error){


		return nil, errors.New("Ssss")
	}


	client:=godis.NewClient("127.0.0.1")
	dd := make(godis.ProtoType)
	dd["dd"]="fff"

	client.Call("hello",&h,dd,shownum)
    waitgroup.Done() //任务完成，将任务队列中的任务数量-1，其实.Done就是.Add(-1)
}




func main (){

	for i := 0; i < 1000; i++ {
			waitgroup.Add(1) //每创建一个goroutine，就把任务队列中任务的数量+1
			go Afunction(i)
		}
		waitgroup.Wait() //.Wait()这里会发生阻塞，直到队列中所有的任务结束就会解除阻塞
     log.Println("finish!!!!")
}