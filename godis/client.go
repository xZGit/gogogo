package godis


import(
	"sync"
	"log"

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



func (c *Client) Call(name string, arg ...interface{})  error {
    c.mutex.Lock()
	defer c.mutex.Unlock()
	c.redisClient.conn.Send("PUBLISH", name, "xxxx")
	c.redisClient.conn.Flush()

	for {

		reply, err := c.redisClient.conn.Receive()
		if err != nil {
			return err
		}
		log.Printf("reply %s", reply)
		// process pushed message
	}



	return nil
}

