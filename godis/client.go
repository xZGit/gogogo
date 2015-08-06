package godis


import(
	"sync"
    "log"
	"time"
	"gopkg.in/redis.v3"
)



// Task handler representation
type taskHandler struct {
	TaskId    string
	HandlerFunc HandleClientFunc
}

type Client struct {
	id string
	redisClient   *RedisClient
	mutex         sync.Mutex
	HandleTasks [] *taskHandler
	hasListen      bool
}



func NewClient(id string, host string) (*Client) {
	client:=Client{
		id:id,
		redisClient:NewRedisClient(host),
		HandleTasks :make([]*taskHandler, 0),
		hasListen:false,
	}

	return &client
}



func (c *Client) Call(name string, handlerFunc HandleClientFunc, args ProtoType,n int)  error {
    c.mutex.Lock()
	defer c.mutex.Unlock()
	if(!c.hasListen){
		go c.Listen()
		c.hasListen = true
	}
	go func() {
//		c.redisClient.getActiveSubPub()
		event, err := newEvent(c.id, args)
		if err != nil {
			panic(err)
		}
		msg, err := event.packBytes()
		task:=taskHandler{
			TaskId:event.MsgId,
			HandlerFunc: handlerFunc,
		}
		c.HandleTasks=append(c.HandleTasks,&task)
     	log.Printf("msg: %v\n", msg)

		c.redisClient.pubConn.Publish(name, string(msg[:]))

	}()
	return nil
}


func (c *Client) Listen() {

	    pubsub, err := c.redisClient.subConn.Subscribe(c.id)
	   if err != nil {
		 panic(err)
	    }

		for {
			msgi, err := pubsub.ReceiveTimeout(100 * time.Millisecond)
			if err != nil {
				err := pubsub.Ping("")
				if err != nil {
					panic(err)
				}
				continue
			}

			switch msg := msgi.(type) {
				case *redis.Subscription:
				log.Println(msg.Kind, msg.Channel)
				case *redis.Message:
//				log.Println(msg.Channel, msg.Payload)
			    go c.ProcessFunc(msg.Payload)
				case *redis.Pong:
//				log.Println(msg)
				default:
				log.Println("unknown message: %#v", msgi)
			}
		}
}



func (c *Client) ProcessFunc(ls string) {

	mySlice := []byte(ls)

	resp, err := unPackRespByte(mySlice)
	if err != nil {
		log.Printf("err: %v\n", err)
	}

	for _, h := range  c.HandleTasks{
		if h.TaskId == resp.MsgId{
			(*h.HandlerFunc)(resp.RespInfo)
		}
	}

}

