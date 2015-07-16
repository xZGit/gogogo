package main


import (
	"./godis"
    "errors"
	"log"
)


func main (){

	exit := make(chan int)



	h := func(v godis.ProtoType) (interface{}, error){
		dd := make(godis.ProtoType)
		dd["dd"]="fff"
		log.Println("yes prcoess v%",v)

		return nil, errors.New("Ssss")
	}


	client:=godis.NewClient("127.0.0.1")
	dd := make(godis.ProtoType)
	dd["dd"]="fff"

	client.Call("hello",&h,dd)
	//    go client.Call("hello")
	<-exit
}