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
	server.RegisterTask("hello", &h)

	server.RegisterTask("hello1", &h)

	client:=godis.NewClient("127.0.0.1")

	go client.Call("hello")

	client.Call("hello")

	<-exit
}