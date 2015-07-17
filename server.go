package main


import (
	"./godis"
)


func main (){

	exit := make(chan int)


	server:=godis.NewServer("127.0.0.1")

	h := func(v godis.ProtoType) (godis.ProtoType, error) {
		dd := make(godis.ProtoType)
		dd["dd"]="fff"
		return dd, nil
	}

	c := make(chan int, 2)
	server.RegisterTask("hello", &h, c)
//	server.RegisterTask("hello1", &h, c)
//	server.RegisterTask("hello2", &h, c)
//	server.RegisterTask("hello3", &h, c)
//	server.RegisterTask("hello4", &h, c)
//	server.RegisterTask("hello5", &h, c)
//	server.RegisterTask("hello6", &h, c)
//	server.RegisterTask("hello7", &h, c)
//	server.RegisterTask("hello8", &h, c)
//	server.RegisterTask("hello9", &h, c)
//	server.RegisterTask("hello10", &h, c)
//	server.RegisterTask("hello11", &h, c)
//	server.RegisterTask("hello12", &h, c)


	<-c

//    go client.Call("hello")
	<-exit
}