package godis


import(
	"sync"
    "log"
	"time"
)


type Client struct {
	redisClient   *RedisClient
	mutex         sync.Mutex

}



func NewClient(host string) (*Client) {
	return  &Client{
		redisClient:NewRedisClient(host),
	}
}



func (c *Client) Call(name string, handlerFunc HandleClientFunc, args ProtoType)  error {
    c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Printf("reply")
	c.redisClient.getActiveSubPub()
	event, err:= newEvent(args)
	if err != nil {
		panic(err)
	}
	rec := make(chan []string)
	go c.redisClient.subConn.Subscribe(rec, event.MsgId)
	msg, err :=event.packBytes()
	c.redisClient.pubConn.Publish(name, msg)
	var ls []string

	for {
		select {
		case ls = <-rec:
			if ls[0]=="message" && len(ls)>2 {
                go c.ProcessFunc(handlerFunc, ls[2])
				go c.redisClient.subConn.Unsubscribe(event.MsgId)
			}
		case <-time.After(5 * time.Second):
//			println("timeout")
			break
		}
	}

	return nil
}



func (c *Client) ProcessFunc(handlerFunc HandleClientFunc, ls string) {

	log.Printf("Client received: %v\n", ls)
	mySlice := []byte(ls)
	log.Printf("Client received: %v\n", mySlice)
	resp, err := unPackRespByte(mySlice)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	if(resp.Code!=0){
		log.Println("return err %s",resp.ErrMsg)
		return
	}

	(*handlerFunc)(resp.Data)

}

