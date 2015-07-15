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


func (c *Client) Call(name string, handlerFunc *func(args []string) (interface{}, error), args ...interface{})  error {
    c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Printf("reply")
	c.redisClient.getActiveSubPub()
	event, err:= newEvent(name, args)
	if err != nil {
		panic(err)
	}

	rec := make(chan []string)
	go c.redisClient.subConn.Subscribe(rec, event.MsgId)
	c.redisClient.pubConn.Publish(name, event)
	var ls []string

	for {
		select {
		case ls = <-rec:
			log.Printf("Client received: %v\n", ls)
		case <-time.After(5 * time.Second):
//			println("timeout")
			break
		}
	}

	return nil
}

