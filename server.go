package main


import (
	"./godis"

)


func main (){

	exit := make(chan int)


	server:=godis.NewServer("127.0.0.1")

	h := func(v []interface{}) (interface{}, error) {

		return "Hello, " + v[0].(string), nil
	}

	c := make(chan int, 1)
	server.RegisterTask("hello", &h, c)

	<-c

	client:=godis.NewClient("127.0.0.1")


	client.Call("hello")
    go client.Call("hello")
	<-exit
}