package main


import (
	"./godis"

)


func main (){

	exit := make(chan int)


	server:=godis.NewServer("127.0.0.1")

	h := func(v []string) (interface{}, error) {

		return "Hello, " + v[0], nil
	}

	c := make(chan int, 2)
	server.RegisterTask("hello", &h, c)
	<-c

//    go client.Call("hello")
	<-exit
}