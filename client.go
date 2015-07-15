package main


import (
	"./godis"

)


func main (){

	exit := make(chan int)



	h := func(v []string) (interface{}, error) {

		return "Hello, " + v[0], nil
	}


	client:=godis.NewClient("127.0.0.1")


	client.Call("hello",&h,"sss","yyyy")
	//    go client.Call("hello")
	<-exit
}