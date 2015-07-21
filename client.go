package main


import (
	"./godis"
    "errors"
	"log"

)

var d=0

func Afunction(client *godis.Client, shownum int) {
	h := func(v godis.RespInfo) (interface{}, error){
		d=d+1
     	log.Println("done: %d",d)
		return nil, errors.New("Ssss")
	}

	dd := make(godis.ProtoType)
	dd["dd"]="fff"

	client.Call("hello",&h,dd,shownum)

}




func main (){
	c := make(chan int)
	client:=godis.NewClient("1", "127.0.0.1")
	for i := 0; i < 3000; i++ {
			go Afunction(client,i)
		}


	log.Println("finish!!!!",<-c)
}