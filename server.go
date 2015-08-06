package main


import (
	"./godis"
)


func main (){

	exit := make(chan int)


	server:=godis.NewServer("127.0.0.1")

	h := func(v godis.ProtoType) (godis.ProtoType, error) {
		result := make(godis.ProtoType)
		result["a"]="some result"
		return result, nil
	}

	c := make(chan int, 2)
	server.RegisterTask("hello", &h, c)
	<-c
	<-exit
}